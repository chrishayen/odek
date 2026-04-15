package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"shotgun.dev/odek/internal/examples"
	"shotgun.dev/odek/internal/toollog"
	"shotgun.dev/odek/openai"
)

type Rune = struct {
	Description   string   `json:"description"`
	FunctionSig   string   `json:"function_signature"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
	Assumptions   []string `json:"assumptions"`
}

type PackageNode struct {
	Name     string          `json:"name"`
	Runes    map[string]Rune `json:"runes"`
	Children []PackageNode   `json:"children,omitempty"`
}

type wirePackage = struct {
	Name  string          `json:"name"`
	Runes map[string]Rune `json:"runes"`
}

type DecompositionResponse struct {
	ProjectPackage wirePackage  `json:"project_package"`
	StdPackage     *wirePackage `json:"std_package,omitempty"`
}

type ClarificationRequest struct {
	Message string `json:"message"`
}

type Client struct {
	api *openai.Client
}

var (
	client          *Client
	decomposeTool   openai.Tool
	readExampleTool openai.Tool
	stdoutMu        sync.Mutex
	exampleIndex    *examples.Index
	exampleManifest string
	toolLogger      *toollog.Logger
)

func init() {
	api, err := openai.NewClient(BASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create openai client: %v\n", err)
		os.Exit(1)
	}
	client = &Client{api: api}

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
			Description: "Submit a rune decomposition. Provide a project_package, and optionally a std_package of reusable utilities.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_package": packageSchema,
					"std_package":     packageSchema,
				},
				"required": []string{"project_package"},
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

	idx, err := examples.LoadFromDir(EXAMPLES_DIR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: could not load example corpus at %s: %v\n", EXAMPLES_DIR, err)
		exampleIndex = &examples.Index{}
	} else {
		exampleIndex = idx
	}
	exampleManifest = exampleIndex.Manifest()

	logger, err := toollog.NewLogger(TOOL_LOG_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: could not open tool log at %s: %v\n", TOOL_LOG_PATH, err)
	} else {
		toolLogger = logger
	}
}

func newConversation(req string) []openai.ChatMessage {
	system := strings.TrimSpace(SYSTEM_PROMPT)
	if exampleManifest != "" {
		system += "\n\n# Example index\n\nThe following reference decompositions are available. To see the full contents of one, call `read_example` with its handle (tier/slug). Pick the most relevant ones for the current requirement, read them first, then call `decompose` with your answer.\n\n" + exampleManifest
	}
	return []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: system},
		{Role: openai.RoleUser, Content: "decompose: " + req},
	}
}

// Decompose drives a multi-turn tool loop with the LLM. The model may call
// `read_example` any number of times (up to MAX_TOOL_ITERATIONS) to retrieve
// reference decompositions from the corpus by handle, then call `decompose`
// to submit its final answer. If the model replies in plain text instead of
// calling a tool, that text is treated as a clarification request.
//
// The full manifest of example handles (tier/slug) is inlined into the
// system prompt at conversation start, so the model picks handles directly
// from context — no search step.
func (c *Client) Decompose(ctx context.Context, messages []openai.ChatMessage) (any, []openai.ChatMessage, error) {
	var parsed *DecompositionResponse

	handler := func(ctx context.Context, call openai.ToolCall) (string, bool, error) {
		switch call.Function.Name {
		case "read_example":
			return handleReadExampleCall(call, messages), false, nil
		case "decompose":
			var d DecompositionResponse
			if err := json.Unmarshal([]byte(call.Function.Arguments), &d); err != nil {
				return "", true, fmt.Errorf("parsing decompose arguments: %w (raw: %s)", err, call.Function.Arguments)
			}
			parsed = &d
			return "decomposition recorded", true, nil
		default:
			return "", true, fmt.Errorf("unexpected tool call: %s", call.Function.Name)
		}
	}

	final, history, err := c.api.AskToolLoop(ctx, messages, []openai.Tool{readExampleTool, decomposeTool}, handler, MAX_TOOL_ITERATIONS)
	if err != nil {
		return nil, history, fmt.Errorf("chat completion failed: %w", err)
	}

	if parsed != nil {
		return *parsed, history, nil
	}

	if len(final.ToolCalls) == 0 && strings.TrimSpace(final.Content) != "" {
		return ClarificationRequest{Message: final.Content}, history, nil
	}
	return nil, history, fmt.Errorf("model returned neither a decompose tool call nor clarification text")
}

func (c *Client) MergeAttempts(ctx context.Context, req string, attempts []DecompositionResponse) (DecompositionResponse, []openai.ChatMessage, error) {
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
		{Role: openai.RoleSystem, Content: strings.TrimSpace(SYSTEM_PROMPT)},
		{Role: openai.RoleUser, Content: userMsg},
	}

	response, history, err := c.Decompose(ctx, messages)
	if err != nil {
		return DecompositionResponse{}, history, err
	}
	decomp, ok := response.(DecompositionResponse)
	if !ok {
		return DecompositionResponse{}, history, fmt.Errorf("merge returned non-decomposition: %T", response)
	}
	return decomp, history, nil
}

func parallelInitialDecompose(ctx context.Context, req string, n int) []DecompositionResponse {
	type attemptResult struct {
		idx  int
		resp DecompositionResponse
		err  error
	}
	out := make(chan attemptResult, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			local := newConversation(req)
			response, _, err := client.Decompose(ctx, local)
			if err != nil {
				out <- attemptResult{i, DecompositionResponse{}, err}
				return
			}
			if decomp, ok := response.(DecompositionResponse); ok && decomp.ProjectPackage.Name != "" {
				out <- attemptResult{i, decomp, nil}
				return
			}
			out <- attemptResult{i, DecompositionResponse{}, fmt.Errorf("non-decomposition response: %T", response)}
		}(i)
	}
	wg.Wait()
	close(out)

	var ok []DecompositionResponse
	for r := range out {
		if r.err != nil {
			stdoutMu.Lock()
			fmt.Printf("   ⚠️  attempt %d failed: %v\n", r.idx+1, r.err)
			stdoutMu.Unlock()
			continue
		}
		ok = append(ok, r.resp)
	}
	return ok
}

type readExampleArgs struct {
	Paths []string `json:"paths"`
}

// handleReadExampleCall parses the tool call, looks up each requested handle
// in the in-memory index, logs the call, and returns the formatted tool
// result (which becomes the Content of the next tool message).
func handleReadExampleCall(call openai.ToolCall, messages []openai.ChatMessage) string {
	var args readExampleArgs
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
		res := exampleIndex.Lookup(ref)
		resolvedList = append(resolvedList, resolved{ref: ref, result: res})
		if res.Entry != nil {
			foundPaths = append(foundPaths, res.Entry.Path)
		}
	}

	if toolLogger != nil {
		_ = toolLogger.LogToolCall(
			time.Now(),
			requirementFromMessages(messages),
			strings.Join(args.Paths, ","),
			len(args.Paths),
			foundPaths,
		)
	}

	stdoutMu.Lock()
	fmt.Printf("🔎 read_example (%d handle%s)\n", len(resolvedList), plural(len(resolvedList)))
	for _, r := range resolvedList {
		switch r.result.Kind {
		case examples.LookupHit:
			fmt.Printf("   ✓ %s\n", r.result.Entry.Handle())
		case examples.LookupTierCorrected:
			fmt.Printf("   ≈ %s (requested %q, corrected to %s)\n",
				r.result.Entry.Handle(), r.ref, r.result.Entry.Handle())
		case examples.LookupMiss:
			fmt.Printf("   ✗ %s (not found)\n", r.ref)
		}
	}
	stdoutMu.Unlock()

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

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
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
