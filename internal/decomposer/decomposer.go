// Package decomposer implements the odek decomposition pipeline as a
// two-pass flow:
//
//  1. Contract pass — the model writes a Design-by-Contract document
//     (purpose, behavior, node hierarchy with +/- tests and assumptions).
//     Streamed back as content deltas.
//  2. Extraction pass — the model calls the `decompose` tool with nested
//     Rune.Children encoding the contract as a typed rune tree.
//
// Callers own a *Decomposer and invoke NewSession to obtain a *Session
// with a live event channel. The producer goroutine runs both passes in
// order; consumers range over Events and call Snapshot() to render.
package decomposer

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"shotgun.dev/odek/internal/examples"
	"shotgun.dev/odek/internal/toollog"
	openai "shotgun.dev/odek/openai"
)

//go:embed contract.md
var contractPromptText string

//go:embed decompose.md
var extractionPromptText string

const maxToolIterations = 6

var (
	decomposeTool   openai.Tool
	readExampleTool openai.Tool
	toolSchemaOnce  sync.Once
)

func initToolSchemas() {
	// Rune schema is recursive (children contains more runes). JSON
	// Schema can express this via $ref + $defs, but OpenAI's tool schema
	// validator is picky — we inline the object twice (top level + one
	// level of children) and let the model handle deeper nesting via the
	// "children" field which we describe but don't structurally type.
	runeLeafSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"description": map[string]any{"type": "string"},
			"function_signature": map[string]any{
				"type":        "string",
				"description": "Bare type signature only, e.g. '(a: i32, b: i32) -> result[i32, string]'. Empty string for parent (non-leaf) runes. Do NOT include any marker prefix like 'fn' or '@'.",
			},
			"positive_tests": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"negative_tests": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"assumptions":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"dependencies":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Fully-qualified paths of runes this rune consumes (e.g. 'std.crypto.hmac_sha256'). Empty when the rune has no dependencies."},
			"children":       map[string]any{"type": "object", "description": "Nested map of child runes keyed by next path segment. Empty object for leaves. Each value has the same shape as this object (recursive)."},
		},
		"required": []string{"description", "function_signature", "positive_tests", "negative_tests", "assumptions", "dependencies", "children"},
	}
	packageSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":  map[string]any{"type": "string"},
			"runes": map[string]any{"type": "object", "additionalProperties": runeLeafSchema},
		},
		"required": []string{"name", "runes"},
	}
	decomposeTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "decompose",
			Description: "Submit a rune decomposition encoding the prior contract. Provide a 1-2 sentence summary, a project_package, and optionally a std_package.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"summary": map[string]any{
						"type":        "string",
						"description": "A 1-2 sentence narrative shown to the user in the chat. Describe what the library does and notable structure. On a refinement, describe what changed and why.",
					},
					"project_package": packageSchema,
					"std_package":     packageSchema,
				},
				"required": []string{"summary", "project_package"},
			},
		},
	}
	readExampleTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "read_example",
			Description: "Read the full contents of one or more example decompositions from the corpus. The full list of available example handles is shown in the `Example index` section of your system message — pick the most relevant ones by name and pass them here. Call this before `decompose` to see how similar requirements have been encoded. You may call it more than once if you need additional references.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"paths": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"minItems":    1,
						"maxItems":    5,
						"description": "list of example handles to load, e.g. ['medium/csv-reader', 'trivial/hello-world']. The handle is tier/slug as shown in the Example index.",
					},
				},
				"required": []string{"paths"},
			},
		},
	}
}

// Decomposer is the shared decomposition client. Construct once with
// NewDecomposer and reuse across NewSession calls.
type Decomposer struct {
	api              *openai.Client
	examples         *examples.Index
	manifest         string
	logger           *toollog.Logger
	contractPrompt   string
	extractionPrompt string
}

// NewDecomposer constructs a Decomposer. Nil-tolerant: passing nil api or
// empty examplesDir / toolLogPath still returns a usable value whose
// NewSession will either succeed (if api is set) or fail with a clear error.
func NewDecomposer(api *openai.Client, examplesDir, toolLogPath string) (*Decomposer, error) {
	toolSchemaOnce.Do(initToolSchemas)

	d := &Decomposer{
		api:              api,
		contractPrompt:   strings.TrimSpace(contractPromptText),
		extractionPrompt: strings.TrimSpace(extractionPromptText),
	}

	if examplesDir != "" {
		idx, err := examples.LoadFromDir(examplesDir)
		if err != nil {
			d.examples = &examples.Index{}
		} else {
			d.examples = idx
		}
	} else {
		d.examples = &examples.Index{}
	}
	d.manifest = d.examples.Manifest()

	if toolLogPath != "" {
		if logger, err := toollog.NewLogger(toolLogPath); err == nil {
			d.logger = logger
		}
	}

	return d, nil
}

// SessionContext is optional refinement context to thread into a /decompose
// call. When either field is set, the conversation the model sees includes
// the prior decomposition and a transcript of the chat discussion.
type SessionContext struct {
	Discussion string
	Prior      *DecompositionResponse
}

// IsEmpty reports whether this SessionContext carries no refinement data.
func (c SessionContext) IsEmpty() bool {
	return c.Discussion == "" && c.Prior == nil
}

// NewSession starts a two-pass decomposition run. It returns a Session
// immediately; the producer goroutine runs pass 1 (contract) then pass 2
// (extraction) in the background, emitting events on sess.Events. Callers
// must range over sess.Events until it closes.
//
// Config is accepted for API stability — the new pipeline produces the
// full tree in one call, so MaxDepth/RuneCap/ParallelInitial are ignored.
func (d *Decomposer) NewSession(ctx context.Context, req string, effortLevel int, effortReason string, cfg Config, sessCtx SessionContext) (*Session, error) {
	_ = cfg
	if d.api == nil {
		return nil, fmt.Errorf("decomposer: no openai client configured")
	}

	sess := newSession(req, effortLevel, effortReason, nil)
	ch := make(chan DecompositionEvent, 32)
	runCtx, cancel := context.WithCancel(ctx)
	sess.setEvents(ch, cancel)

	go d.run(runCtx, sess, sessCtx, ch)

	return sess, nil
}

// run is the producer goroutine. It drives pass 1, then pass 2, emitting
// events onto ch. Emits EventDone before closing the channel.
func (d *Decomposer) run(ctx context.Context, sess *Session, sessCtx SessionContext, ch chan<- DecompositionEvent) {
	start := time.Now()
	defer func() {
		emit := func(evt DecompositionEvent) {
			sess.Apply(evt)
			select {
			case ch <- evt:
			default:
			}
		}
		emit(EventDone{ElapsedMs: time.Since(start).Milliseconds()})
		close(ch)
		sess.clearEvents()
	}()

	emit := func(evt DecompositionEvent) {
		sess.Apply(evt)
		select {
		case ch <- evt:
		case <-ctx.Done():
		}
	}

	// Pass 1: contract.
	emit(EventPhaseStarted{Phase: PhaseContract})
	contract, contractMsgs, err := d.runContract(ctx, sess.Requirement, sessCtx, emit)
	if err != nil {
		if ctx.Err() != nil {
			emit(EventCancelled{})
			return
		}
		emit(EventError{Phase: PhaseContract, Err: err.Error()})
		return
	}
	sess.BaseMessages = contractMsgs

	// Pass 2: extraction.
	emit(EventPhaseStarted{Phase: PhaseExtraction})
	resp, err := d.runExtraction(ctx, sess.Requirement, contract, sessCtx, emit)
	if err != nil {
		if ctx.Err() != nil {
			emit(EventCancelled{})
			return
		}
		emit(EventError{Phase: PhaseExtraction, Err: err.Error()})
		return
	}
	emit(EventRunesComplete{Response: resp, ElapsedMs: time.Since(start).Milliseconds()})
}

// runContract executes pass 1: a streaming content call with the contract
// system prompt. Emits EventContractChunk on each delta, EventContractComplete
// at the end. Returns the full contract text and the conversation history.
func (d *Decomposer) runContract(ctx context.Context, req string, sessCtx SessionContext, emit func(DecompositionEvent)) (string, []openai.ChatMessage, error) {
	userMsg := buildContractUserMessage(req, sessCtx)
	msgs := []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: d.contractPrompt},
		{Role: openai.RoleUser, Content: userMsg},
	}

	start := time.Now()
	chunkCtx := openai.WithContentCallback(ctx, func(delta string) {
		emit(EventContractChunk{Text: delta})
	})

	resp, err := d.api.Chat(chunkCtx, &openai.ChatCompletionRequest{
		Model:    openai.DefaultModel,
		Messages: msgs,
	})
	if err != nil {
		return "", msgs, fmt.Errorf("contract pass: %w", err)
	}
	full := strings.TrimSpace(resp.Choices[0].Message.Content)
	if full == "" {
		return "", msgs, fmt.Errorf("contract pass returned empty content")
	}

	emit(EventContractComplete{Full: full, ElapsedMs: time.Since(start).Milliseconds()})

	history := append(msgs, openai.ChatMessage{
		Role:    openai.RoleAssistant,
		Content: full,
	})
	return full, history, nil
}

// buildContractUserMessage assembles the user message for pass 1. On a
// refinement pass it includes the prior decomposition (as JSON) and the
// chat discussion so the model can revise.
func buildContractUserMessage(req string, sessCtx SessionContext) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Requirement: %s\n\n", req)
	if sessCtx.Prior != nil {
		priorJSON, err := json.MarshalIndent(sessCtx.Prior, "", "  ")
		if err == nil {
			fmt.Fprintf(&b, "Prior decomposition (encoded form of a previous contract):\n```json\n%s\n```\n\n", string(priorJSON))
		}
	}
	if sessCtx.Discussion != "" {
		label := "Conversation with the user leading up to this requirement:"
		if sessCtx.Prior != nil {
			label = "Discussion with the user since the prior decomposition:"
		}
		fmt.Fprintf(&b, "%s\n%s\n\n", label, strings.TrimSpace(sessCtx.Discussion))
	}
	if sessCtx.Prior != nil {
		b.WriteString("This is a refinement pass. Revise the prior decomposition's contract according to the discussion above: preserve what still applies, restructure or rename where the discussion asks for it, and add or remove nodes as needed. Produce the updated contract.")
	} else {
		b.WriteString("Write the contract now.")
	}
	return b.String()
}

// runExtraction executes pass 2: a tool-loop call with the extraction
// system prompt, with the contract included as input. Streams tool-arg
// bytes via EventExtractionProgress. Returns the parsed DecompositionResponse.
func (d *Decomposer) runExtraction(ctx context.Context, req string, contract string, sessCtx SessionContext, emit func(DecompositionEvent)) (*DecompositionResponse, error) {
	system := d.extractionPrompt
	if d.manifest != "" {
		system += "\n\n# Example index\n\nThe following reference decompositions are available for stylistic anchoring. Call `read_example` with tier/slug handles if you want to see how similar trees have been shaped. Only use it if you need style guidance — the contract is authoritative on structure.\n\n" + d.manifest
	}

	userMsg := fmt.Sprintf("Requirement: %s\n\nContract:\n%s\n\nEncode this contract by calling the `decompose` tool.", req, contract)
	_ = sessCtx // refinement context was consumed in pass 1; pass 2 just encodes the produced contract
	msgs := []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: system},
		{Role: openai.RoleUser, Content: userMsg},
	}

	var parsed *DecompositionResponse
	var argBytes int

	argsCtx := openai.WithToolArgsCallback(ctx, func(_ int, delta string) {
		argBytes += len(delta)
		emit(EventExtractionProgress{Bytes: argBytes})
	})

	handler := func(ctx context.Context, call openai.ToolCall) (string, bool, error) {
		switch call.Function.Name {
		case "read_example":
			return d.handleReadExampleCall(call, req, emit), false, nil
		case "decompose":
			var dr DecompositionResponse
			if err := json.Unmarshal([]byte(call.Function.Arguments), &dr); err != nil {
				return "", true, fmt.Errorf("parsing decompose arguments: %w (raw: %s)", err, call.Function.Arguments)
			}
			normalizePackageSignatures(&dr.ProjectPackage)
			if dr.StdPackage != nil {
				normalizePackageSignatures(dr.StdPackage)
			}
			parsed = &dr
			return "decomposition recorded", true, nil
		default:
			return "", true, fmt.Errorf("unexpected tool call: %s", call.Function.Name)
		}
	}

	tools := []openai.Tool{readExampleTool, decomposeTool}
	final, _, err := d.api.AskToolLoop(argsCtx, msgs, tools, handler, maxToolIterations, nil)
	if err != nil {
		return nil, fmt.Errorf("extraction pass: %w", err)
	}
	if parsed != nil {
		return parsed, nil
	}

	// Plain-text reply without a tool call: force a retry with required tool use.
	if len(final.ToolCalls) == 0 && strings.TrimSpace(final.Content) != "" {
		msgs = append(msgs, openai.ChatMessage{
			Role:    openai.RoleAssistant,
			Content: final.Content,
		}, openai.ChatMessage{
			Role: openai.RoleUser,
			Content: "Your previous response was plain text. You MUST submit your answer " +
				"by calling the `decompose` tool. Encode the contract above as the tool's " +
				"nested runes tree. Call `decompose` now.",
		})
		_, _, retryErr := d.api.AskToolLoop(argsCtx, msgs, tools, handler, maxToolIterations, "required")
		if retryErr != nil {
			return nil, fmt.Errorf("forced-retry after plain-text reply failed: %w", retryErr)
		}
		if parsed != nil {
			return parsed, nil
		}
	}
	return nil, fmt.Errorf("model returned no decompose tool call")
}

// handleReadExampleCall parses a read_example tool call, resolves each
// handle, logs the call, emits EventReadExample via emit, and returns the
// formatted tool result.
func (d *Decomposer) handleReadExampleCall(call openai.ToolCall, requirement string, emit func(DecompositionEvent)) string {
	var args struct {
		Paths []string `json:"paths"`
	}
	if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
		return fmt.Sprintf("error: could not parse read_example arguments: %v", err)
	}
	if len(args.Paths) == 0 {
		return "error: `paths` is required and must contain at least one entry"
	}
	if len(args.Paths) > 5 {
		args.Paths = args.Paths[:5]
	}

	type resolved struct {
		ref    string
		result examples.LookupResult
	}
	resolvedList := make([]resolved, 0, len(args.Paths))
	foundPaths := make([]string, 0, len(args.Paths))
	for _, ref := range args.Paths {
		res := d.examples.Lookup(ref)
		resolvedList = append(resolvedList, resolved{ref: ref, result: res})
		if res.Entry != nil {
			foundPaths = append(foundPaths, res.Entry.Path)
		}
	}

	if d.logger != nil {
		_ = d.logger.LogToolCall(
			time.Now(),
			requirement,
			strings.Join(args.Paths, ","),
			len(args.Paths),
			foundPaths,
		)
	}
	if emit != nil {
		emit(EventReadExample{Paths: args.Paths, Found: foundPaths})
	}

	var b strings.Builder
	for i, r := range resolvedList {
		switch r.result.Kind {
		case examples.LookupHit:
			fmt.Fprintf(&b, "=== %s (tier=%s) ===\n", r.result.Entry.Handle(), r.result.Entry.Tier)
			b.WriteString(r.result.Entry.Content)
			if !strings.HasSuffix(r.result.Entry.Content, "\n") {
				b.WriteString("\n")
			}
			b.WriteString("\n")
		case examples.LookupTierCorrected:
			fmt.Fprintf(&b, "=== %s (tier=%s, auto-corrected from %q — the slug lives in a different tier than you guessed) ===\n",
				r.result.Entry.Handle(), r.result.Entry.Tier, r.ref)
			b.WriteString(r.result.Entry.Content)
			if !strings.HasSuffix(r.result.Entry.Content, "\n") {
				b.WriteString("\n")
			}
			b.WriteString("\n")
		case examples.LookupMiss:
			fmt.Fprintf(&b, "=== request %d: %q NOT FOUND in example index ===\n", i+1, r.ref)
			if len(r.result.Suggestions) > 0 {
				b.WriteString("Did you mean one of these?\n")
				for _, s := range r.result.Suggestions {
					fmt.Fprintf(&b, "  - %s\n", s.Handle())
				}
			} else {
				b.WriteString("(no similar handles found; try a different slug from the Example index)\n")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}
