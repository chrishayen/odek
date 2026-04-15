// Package decomposer implements the odek decomposition pipeline: tool-calling
// client that talks to the model via the `decompose` and `read_example`
// tools, plus recursive expansion of the resulting rune tree. Callers own a
// *Decomposer and invoke NewSession to obtain a *Session; the same Session
// flows through the TUI (chat + decomposition page) and the CLI.
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

//go:embed decompose.md
var systemPromptText string

const maxToolIterations = 6

var (
	decomposeTool   openai.Tool
	readExampleTool openai.Tool
	toolSchemaOnce  sync.Once
)

func initToolSchemas() {
	runeSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"description":        map[string]any{"type": "string"},
			"function_signature": map[string]any{"type": "string"},
			"positive_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"negative_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"assumptions":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		},
		"required": []string{"description", "function_signature", "positive_tests", "negative_tests", "assumptions"},
	}
	packageSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":  map[string]any{"type": "string"},
			"runes": map[string]any{"type": "object", "additionalProperties": runeSchema},
		},
		"required": []string{"name", "runes"},
	}
	decomposeTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "decompose",
			Description: "Submit a rune decomposition. Provide a 1-2 sentence summary, a project_package, and optionally a std_package of reusable utilities.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"summary": map[string]any{
						"type":        "string",
						"description": "A 1-2 sentence narrative shown to the user in the chat. On a fresh decomposition, describe what the feature is and the approach you took. On a refinement pass (when a prior decomposition is included), describe what you changed in response to the user's latest feedback and why. Explain, do not list rune names.",
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
			Description: "Read the full contents of one or more example decompositions from the corpus. The full list of available example handles is shown in the `Example index` section of your system message — pick the most relevant ones by name and pass them here. Call this before `decompose` to see how similar requirements have been broken down. You may call it more than once if you need additional references.",
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
// NewDecomposer and reuse across decompose/expand calls.
type Decomposer struct {
	api          *openai.Client
	examples     *examples.Index
	manifest     string
	logger       *toollog.Logger
	systemPrompt string
}

// NewDecomposer constructs a Decomposer. Nil-tolerant: passing nil api
// or empty examplesDir / toolLogPath still returns a usable value whose
// Decompose will either succeed (if api is set) or fail with a clear error.
// This is useful for unit tests that render without running the pipeline.
func NewDecomposer(api *openai.Client, examplesDir, toolLogPath string) (*Decomposer, error) {
	toolSchemaOnce.Do(initToolSchemas)

	d := &Decomposer{
		api:          api,
		systemPrompt: strings.TrimSpace(systemPromptText),
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
// the prior decomposition (as JSON) and a transcript of the chat discussion
// alongside the original requirement, so the model refines the prior answer
// instead of starting from scratch.
type SessionContext struct {
	// Discussion is a pre-formatted transcript of chat turns since the
	// original requirement was stated, e.g.
	//   you: add a JSON section
	//   clank: okay, added json handling...
	// The caller owns the formatting.
	Discussion string

	// Prior is the DecompositionResponse from the previous /decompose run,
	// included verbatim as the starting point for refinement.
	Prior *DecompositionResponse
}

// IsEmpty reports whether this SessionContext carries no refinement data.
func (c SessionContext) IsEmpty() bool {
	return c.Discussion == "" && c.Prior == nil
}

// newConversation builds the initial system+user messages for a fresh
// decompose call. When sessCtx has discussion or a prior decomposition,
// the user message is a refinement prompt that embeds both.
func (d *Decomposer) newConversation(req string, sessCtx SessionContext) []openai.ChatMessage {
	system := d.systemPrompt
	if d.manifest != "" {
		system += "\n\n# Example index\n\nThe following reference decompositions are available. To see the full contents of one, call `read_example` with its handle (tier/slug). Pick the most relevant ones for the current requirement, read them first, then call `decompose` with your answer.\n\n" + d.manifest
	}

	userContent := "decompose: " + req
	if !sessCtx.IsEmpty() {
		userContent = buildRefinementMessage(req, sessCtx)
	}

	return []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: system},
		{Role: openai.RoleUser, Content: userContent},
	}
}

// buildRefinementMessage constructs the user message for a refinement pass.
// Layout: original requirement, prior decomposition (as JSON), chat
// discussion, then an instruction that adapts to whether a prior
// decomposition is present.
func buildRefinementMessage(req string, sessCtx SessionContext) string {
	var b strings.Builder
	fmt.Fprintf(&b, "decompose: %s\n\n", req)
	if sessCtx.Prior != nil {
		priorJSON, err := json.MarshalIndent(sessCtx.Prior, "", "  ")
		if err == nil {
			fmt.Fprintf(&b, "Prior decomposition:\n```json\n%s\n```\n\n", string(priorJSON))
		}
	}
	if sessCtx.Discussion != "" {
		label := "Conversation with the user leading up to this request:"
		if sessCtx.Prior != nil {
			label = "Discussion with the user since the prior decomposition:"
		}
		fmt.Fprintf(&b, "%s\n%s\n\n", label, strings.TrimSpace(sessCtx.Discussion))
	}
	if sessCtx.Prior != nil {
		b.WriteString("This is a refinement pass. The prior decomposition is the starting point; the discussion describes what the user wants changed. Preserve the parts that still apply, rename or restructure where the discussion asks for it, and add or remove runes as needed. Submit the refined decomposition via the decompose tool.")
	} else {
		b.WriteString("The conversation above is context for this decomposition. Use it to inform scope, naming, and assumptions, then submit the decomposition via the decompose tool.")
	}
	return b.String()
}

// Decompose drives a multi-turn tool loop. The model may call read_example
// any number of times (up to maxToolIterations), then calls the decompose
// tool with its final answer. If it replies in plain text instead, that
// text is returned as a ClarificationRequest. Returns the parsed
// DecompositionResponse (as any), the full message history, and any error.
//
// emit is called for every read_example tool call. Pass nil for no-op.
func (d *Decomposer) Decompose(ctx context.Context, messages []openai.ChatMessage, emit func(ExpansionEvent)) (any, []openai.ChatMessage, error) {
	if d.api == nil {
		return nil, nil, fmt.Errorf("decomposer: no openai client configured")
	}

	var parsed *DecompositionResponse
	if emit == nil {
		emit = func(ExpansionEvent) {}
	}

	handler := func(ctx context.Context, call openai.ToolCall) (string, bool, error) {
		switch call.Function.Name {
		case "read_example":
			result := d.handleReadExampleCall(call, messages, emit)
			return result, false, nil
		case "decompose":
			var dr DecompositionResponse
			if err := json.Unmarshal([]byte(call.Function.Arguments), &dr); err != nil {
				return "", true, fmt.Errorf("parsing decompose arguments: %w (raw: %s)", err, call.Function.Arguments)
			}
			parsed = &dr
			return "decomposition recorded", true, nil
		default:
			return "", true, fmt.Errorf("unexpected tool call: %s", call.Function.Name)
		}
	}

	tools := []openai.Tool{readExampleTool, decomposeTool}
	final, history, err := d.api.AskToolLoop(ctx, messages, tools, handler, maxToolIterations, nil)
	if err != nil {
		return nil, history, fmt.Errorf("chat completion failed: %w", err)
	}

	if parsed != nil {
		return *parsed, history, nil
	}

	// Plain-text reply without a tool call: the model ignored the decompose
	// tool schema. Local llama.cpp models pattern-match the prompt's tree
	// examples and emit trees in prose. Append a corrective user message and
	// force a tool call on a second pass; if that still fails, fall through
	// to ClarificationRequest.
	if len(final.ToolCalls) == 0 && strings.TrimSpace(final.Content) != "" {
		history = append(history, openai.ChatMessage{
			Role: openai.RoleUser,
			Content: "Your previous response was plain text. You MUST submit your answer " +
				"by calling the `decompose` tool. Convert the structure you just described " +
				"into the JSON arguments the tool expects: a `summary` string (1-2 sentence " +
				"narrative), a `project_package` object, and optionally a `std_package` " +
				"object. Each package has a `name` and a `runes` map; every rune has " +
				"`description`, `function_signature`, `positive_tests`, `negative_tests`, " +
				"and `assumptions` fields. Call `decompose` now.",
		})
		retryFinal, retryHistory, retryErr := d.api.AskToolLoop(ctx, history, tools, handler, maxToolIterations, "required")
		history = retryHistory
		if retryErr != nil {
			return nil, history, fmt.Errorf("forced-retry after plain-text reply failed: %w", retryErr)
		}
		if parsed != nil && parsed.ProjectPackage.Name != "" && len(parsed.ProjectPackage.Runes) > 0 {
			return *parsed, history, nil
		}
		// Retry produced either more plain text or degenerate JSON. Fall
		// through to the clarification path using whichever final turn is
		// useful.
		if len(retryFinal.ToolCalls) == 0 && strings.TrimSpace(retryFinal.Content) != "" {
			return ClarificationRequest{Message: retryFinal.Content}, history, nil
		}
		return ClarificationRequest{Message: final.Content}, history, nil
	}
	return nil, history, fmt.Errorf("model returned neither a decompose tool call nor clarification text")
}

// MergeAttempts asks the model to merge N independent decompositions into a
// single consensus. Used when parallel initial attempts come back with
// slightly different structures.
func (d *Decomposer) MergeAttempts(ctx context.Context, req string, attempts []DecompositionResponse) (DecompositionResponse, []openai.ChatMessage, error) {
	var blocks []string
	for i, a := range attempts {
		b, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			return DecompositionResponse{}, nil, err
		}
		blocks = append(blocks, fmt.Sprintf("Attempt %d:\n%s", i+1, string(b)))
	}

	userMsg := fmt.Sprintf(`Below are %d independent decompositions of this requirement:

REQUIREMENT: %s

Merge them into a single consensus decomposition. Take the best ideas from each, drop redundancy, prefer the clearest names. The result should be a single project_package (and optional std_package) that captures the agreed-on top-level architecture.

Submit the consensus by calling the decompose tool.

%s`, len(attempts), req, strings.Join(blocks, "\n\n"))

	messages := []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: d.systemPrompt},
		{Role: openai.RoleUser, Content: userMsg},
	}

	response, history, err := d.Decompose(ctx, messages, nil)
	if err != nil {
		return DecompositionResponse{}, history, err
	}
	decomp, ok := response.(DecompositionResponse)
	if !ok {
		return DecompositionResponse{}, history, fmt.Errorf("merge returned non-decomposition: %T", response)
	}
	return decomp, history, nil
}

// parallelInitialDecompose runs n concurrent initial Decompose calls. It
// returns the successful decompositions and any clarification messages
// the model produced along the way, so NewSession can fall back to a
// clarification if every attempt declined to decompose.
func (d *Decomposer) parallelInitialDecompose(ctx context.Context, req string, n int, sessCtx SessionContext) (successes []DecompositionResponse, clarifications []string) {
	type attemptResult struct {
		idx           int
		resp          DecompositionResponse
		clarification string
		err           error
	}
	out := make(chan attemptResult, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			local := d.newConversation(req, sessCtx)
			response, _, err := d.Decompose(ctx, local, nil)
			if err != nil {
				out <- attemptResult{idx: i, err: err}
				return
			}
			if clar, ok := response.(ClarificationRequest); ok {
				out <- attemptResult{idx: i, clarification: clar.Message}
				return
			}
			if decomp, ok := response.(DecompositionResponse); ok && decomp.ProjectPackage.Name != "" {
				out <- attemptResult{idx: i, resp: decomp}
				return
			}
			out <- attemptResult{idx: i, err: fmt.Errorf("non-decomposition response: %T", response)}
		}(i)
	}
	wg.Wait()
	close(out)

	for r := range out {
		switch {
		case r.clarification != "":
			clarifications = append(clarifications, r.clarification)
		case r.err == nil:
			successes = append(successes, r.resp)
		}
	}
	return successes, clarifications
}

// NewSession runs the initial decomposition for a requirement and returns a
// Session wrapping the root. Does NOT kick off recursion — callers must
// invoke ExpandStreaming separately when they want deeper levels.
//
// When sessCtx carries refinement context (discussion and/or a prior
// decomposition), the model sees that context in the user message and
// treats this call as a refinement pass rather than a fresh decompose.
// Passing a zero-value SessionContext is equivalent to the first-time
// decompose behavior.
func (d *Decomposer) NewSession(ctx context.Context, req string, effortLevel int, effortReason string, cfg Config, sessCtx SessionContext) (*Session, error) {
	if d.api == nil {
		return nil, fmt.Errorf("decomposer: no openai client configured")
	}

	var root DecompositionResponse
	var baseMessages []openai.ChatMessage

	if cfg.ParallelInitial <= 1 {
		baseMessages = d.newConversation(req, sessCtx)
		response, history, err := d.Decompose(ctx, baseMessages, nil)
		if err != nil {
			return nil, fmt.Errorf("initial decompose: %w", err)
		}
		if clar, isClar := response.(ClarificationRequest); isClar {
			return nil, &ClarificationNeeded{Message: clar.Message}
		}
		decomp, ok := response.(DecompositionResponse)
		if !ok {
			return nil, fmt.Errorf("unexpected response type %T", response)
		}
		root = decomp
		baseMessages = history
	} else {
		attempts, clarifications := d.parallelInitialDecompose(ctx, req, cfg.ParallelInitial, sessCtx)
		if len(attempts) == 0 {
			if len(clarifications) > 0 {
				return nil, &ClarificationNeeded{Message: clarifications[0]}
			}
			return nil, fmt.Errorf("all parallel attempts failed")
		}
		if len(attempts) == 1 {
			root = attempts[0]
			baseMessages = d.newConversation(req, sessCtx)
		} else {
			merged, mergedMsgs, err := d.MergeAttempts(ctx, req, attempts)
			if err != nil {
				root = attempts[0]
				baseMessages = d.newConversation(req, sessCtx)
			} else {
				root = merged
				baseMessages = mergedMsgs
			}
		}
	}

	rootAuto := &AutoDecomposition{
		Path:       "root",
		Depth:      0,
		Response:   &root,
		ParentPath: "",
		ChildPaths: make([]string, 0),
	}
	return newSession(req, effortLevel, effortReason, rootAuto, baseMessages), nil
}

// handleReadExampleCall parses the tool call, resolves each handle, logs
// the call, emits an EventReadExample via the provided emitter, and returns
// the formatted tool result (which becomes the Content of the next tool
// message).
func (d *Decomposer) handleReadExampleCall(call openai.ToolCall, messages []openai.ChatMessage, emit func(ExpansionEvent)) string {
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
			requirementFromMessages(messages),
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

// requirementFromMessages walks the conversation backward and returns the
// first user message content that starts with "decompose: " (the convention
// from newConversation). Falls back to the most recent user message text.
func requirementFromMessages(messages []openai.ChatMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		if m.Role != openai.RoleUser {
			continue
		}
		content := strings.TrimSpace(m.Content)
		if after, ok := strings.CutPrefix(content, "decompose: "); ok {
			return after
		}
	}
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		if m.Role == openai.RoleUser {
			return strings.TrimSpace(m.Content)
		}
	}
	return ""
}
