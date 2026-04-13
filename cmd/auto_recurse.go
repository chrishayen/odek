package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	BASE_URL   = "http://localhost:8080/v1"
	MODEL_NAME = "default"
)

//go:embed decompose.md
var SYSTEM_PROMPT string

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

type EffortEstimate struct {
	Level  int    `json:"level"`
	Reason string `json:"reason"`
}

type RunConfig struct {
	ParallelInitial int
	MaxDepth        int
	RuneCap         int
	Recurse         bool
}

type Client struct {
	openai *openai.Client
}

type RuneExpansionInfo struct {
	FullPath            string
	Depth               int
	ParentDecomposition *AutoDecomposition
}

type AutoDecomposition struct {
	Path       string
	Depth      int
	Response   *DecompositionResponse
	ParentPath string
	ChildPaths []string
}

var (
	client         *Client
	decomposeTool  openai.Tool
	rateEffortTool openai.Tool
	stdoutMu       sync.Mutex
)

func init() {
	config := openai.DefaultConfig("default")
	config.BaseURL = BASE_URL
	client = &Client{openai: openai.NewClientWithConfig(config)}

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

	rateEffortTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "rate_effort",
			Description: "Rate the complexity of a software requirement on a 1-5 scale.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"level": map[string]any{
						"type":        "integer",
						"minimum":     1,
						"maximum":     5,
						"description": "1=trivial (hello world, single function); 2=small (one file or simple CLI); 3=medium (a few modules); 4=large (subsystem with several integration points); 5=very large (full application stack)",
					},
					"reason": map[string]any{
						"type":        "string",
						"description": "One short sentence justifying the level.",
					},
				},
				"required": []string{"level", "reason"},
			},
		},
	}
}

func getRequirement(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)
	return requirement, nil
}

func main() {
	ctx := context.Background()
	printBanner()

	stdin := bufio.NewReader(os.Stdin)
	req, err := getRequirement(stdin)
	if err != nil || req == "" {
		fmt.Println("No requirement provided. Exiting.")
		os.Exit(1)
	}

	// Stage 1: estimate effort
	fmt.Printf("\n🎯 Estimating effort for: %s\n", req)
	effortStart := time.Now()
	effort, err := client.EstimateEffort(ctx, req)
	if err != nil {
		fmt.Printf("⚠️  effort estimation failed: %v — defaulting to level 3\n", err)
		effort = EffortEstimate{Level: 3, Reason: "default (estimator failed)"}
	}
	cfg := configForEffort(effort.Level)
	fmt.Printf("🎯 effort: %d/5 — %s (%s)\n", effort.Level, effort.Reason, time.Since(effortStart).Round(time.Millisecond))
	fmt.Printf("   config: parallel=%d depth=%d cap=%d recurse=%v\n",
		cfg.ParallelInitial, cfg.MaxDepth, cfg.RuneCap, cfg.Recurse)

	// Stages 2 & 3: initial decompose (single or N-way + merge)
	var rootResponse DecompositionResponse
	var baseMessages []openai.ChatCompletionMessage

	if cfg.ParallelInitial == 1 {
		baseMessages = newConversation(req)
		fmt.Printf("\nDecomposing: %s...\n", req)
		initStart := time.Now()
		response, err := client.Decompose(ctx, &baseMessages)
		fmt.Printf("⏱️  initial decompose: %s\n", time.Since(initStart).Round(time.Millisecond))
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
		switch v := response.(type) {
		case ClarificationRequest:
			fmt.Printf("\n⚠️  CLARIFICATION NEEDED: %s\n", v.Message)
			return
		case DecompositionResponse:
			rootResponse = v
		default:
			fmt.Printf("ERROR: unexpected response type %T\n", response)
			return
		}
	} else {
		fmt.Printf("\n🚀 Running %d parallel initial decompositions...\n", cfg.ParallelInitial)
		parallelStart := time.Now()
		attempts := parallelInitialDecompose(ctx, req, cfg.ParallelInitial)
		fmt.Printf("✅ %d/%d attempts succeeded (%s)\n",
			len(attempts), cfg.ParallelInitial, time.Since(parallelStart).Round(time.Millisecond))
		if len(attempts) == 0 {
			fmt.Println("ERROR: all parallel attempts failed")
			return
		}
		if len(attempts) == 1 {
			rootResponse = attempts[0]
			baseMessages = newConversation(req)
		} else {
			fmt.Printf("\n🔀 Merging %d attempts...\n", len(attempts))
			mergeStart := time.Now()
			merged, mergedMsgs, err := client.MergeAttempts(ctx, req, attempts)
			fmt.Printf("⏱️  merge: %s\n", time.Since(mergeStart).Round(time.Millisecond))
			if err != nil {
				fmt.Printf("⚠️  merge failed: %v — using first attempt\n", err)
				rootResponse = attempts[0]
				baseMessages = newConversation(req)
			} else {
				rootResponse = merged
				baseMessages = mergedMsgs
			}
		}
	}

	// Stage 4: present + confirm
	root := &AutoDecomposition{
		Path:       "root",
		Depth:      0,
		Response:   &rootResponse,
		ParentPath: "",
		ChildPaths: make([]string, 0),
	}
	queue := collectRunesForExpansion(root.Response)

	printInitialDecomposition(root.Response)

	if !cfg.Recurse {
		fmt.Println("\n(Skipping recursion: requirement is trivial enough.)")
		return
	}
	if len(queue) == 0 {
		fmt.Println("\n(No runes to expand.)")
		return
	}

	if !confirm(stdin, fmt.Sprintf("\nProceed with recursion (max depth %d, max %d runes)? [y/N] ", cfg.MaxDepth, cfg.RuneCap)) {
		fmt.Println("Stopping before recursion. Initial decomposition is above.")
		return
	}

	// Stage 5: recurse
	for i := range queue {
		queue[i].ParentDecomposition = root
	}

	expandRecursively(ctx, baseMessages, root, queue, cfg)
}

func configForEffort(level int) RunConfig {
	switch level {
	case 1:
		return RunConfig{ParallelInitial: 1, MaxDepth: 0, RuneCap: 10, Recurse: false}
	case 2:
		return RunConfig{ParallelInitial: 1, MaxDepth: 1, RuneCap: 25, Recurse: true}
	case 3:
		return RunConfig{ParallelInitial: 3, MaxDepth: 2, RuneCap: 50, Recurse: true}
	case 4:
		return RunConfig{ParallelInitial: 5, MaxDepth: 3, RuneCap: 100, Recurse: true}
	case 5:
		return RunConfig{ParallelInitial: 5, MaxDepth: 3, RuneCap: 200, Recurse: true}
	}
	return RunConfig{ParallelInitial: 3, MaxDepth: 2, RuneCap: 50, Recurse: true}
}

func confirm(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

func parallelInitialDecompose(ctx context.Context, req string, n int) []DecompositionResponse {
	type attemptResult struct {
		idx  int
		resp DecompositionResponse
		err  error
	}
	out := make(chan attemptResult, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			local := newConversation(req)
			response, err := client.Decompose(ctx, &local)
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

func (c *Client) EstimateEffort(ctx context.Context, req string) (EffortEstimate, error) {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: "You are a software-complexity estimator. Given a software requirement, rate it 1-5 by calling the rate_effort tool. Reply only via the tool call."},
		{Role: openai.ChatMessageRoleUser, Content: "Rate the complexity of this requirement: " + req},
	}
	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    MODEL_NAME,
		Messages: messages,
		Tools:    []openai.Tool{rateEffortTool},
		ToolChoice: map[string]any{
			"type": "function",
			"function": map[string]any{
				"name": "rate_effort",
			},
		},
	})
	if err != nil {
		return EffortEstimate{}, fmt.Errorf("effort completion failed: %w", err)
	}
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		return EffortEstimate{}, fmt.Errorf("no tool call in effort response")
	}
	call := resp.Choices[0].Message.ToolCalls[0]
	var est EffortEstimate
	if err := json.Unmarshal([]byte(call.Function.Arguments), &est); err != nil {
		return EffortEstimate{}, fmt.Errorf("parsing effort args: %w (raw: %s)", err, call.Function.Arguments)
	}
	if est.Level < 1 || est.Level > 5 {
		return EffortEstimate{}, fmt.Errorf("level out of range: %d", est.Level)
	}
	return est, nil
}

func (c *Client) MergeAttempts(ctx context.Context, req string, attempts []DecompositionResponse) (DecompositionResponse, []openai.ChatCompletionMessage, error) {
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

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: strings.TrimSpace(SYSTEM_PROMPT)},
		{Role: openai.ChatMessageRoleUser, Content: userMsg},
	}

	response, err := c.Decompose(ctx, &messages)
	if err != nil {
		return DecompositionResponse{}, nil, err
	}
	decomp, ok := response.(DecompositionResponse)
	if !ok {
		return DecompositionResponse{}, nil, fmt.Errorf("merge returned non-decomposition: %T", response)
	}
	return decomp, messages, nil
}

func printInitialDecomposition(resp *DecompositionResponse) {
	fmt.Printf("\n🌳 INITIAL DECOMPOSITION:\n")
	if len(resp.ProjectPackage.Runes) > 0 {
		fmt.Printf("   📦 %s\n", resp.ProjectPackage.Name)
		printRunesIndented(resp.ProjectPackage.Runes, 1)
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		fmt.Printf("   📚 %s\n", resp.StdPackage.Name)
		printRunesIndented(resp.StdPackage.Runes, 1)
	}
}

func printBanner() {
	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
}

func newConversation(req string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: strings.TrimSpace(SYSTEM_PROMPT)},
		{Role: openai.ChatMessageRoleUser, Content: "decompose: " + req},
	}
}

type expansionResult struct {
	runeInfo RuneExpansionInfo
	resp     *DecompositionResponse
	err      error
}

// expandRecursively drains the expansion queue level-by-level, decomposing each
// rune up to cfg.MaxDepth and stitching child decompositions back into the tree.
// Each level dispatches all expansions in parallel since each expansion is
// independent — it only needs the initial decomposition context, not the
// results of sibling expansions.
func expandRecursively(ctx context.Context, baseMessages []openai.ChatCompletionMessage, root *AutoDecomposition, queue []RuneExpansionInfo, cfg RunConfig) {
	for i := range queue {
		queue[i].ParentDecomposition = root
	}

	allDecompositions := []*AutoDecomposition{root}
	totalRunesCount := countTotalRunes(root.Response)
	visitedRunePaths := map[string]bool{"root": true}

	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(queue))

	currentLevel := queue
	for len(currentLevel) > 0 {
		if totalRunesCount >= cfg.RuneCap {
			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", cfg.RuneCap)
			break
		}

		var toExpand []RuneExpansionInfo
		for _, ri := range currentLevel {
			if visitedRunePaths[ri.FullPath] || ri.Depth >= cfg.MaxDepth {
				continue
			}
			visitedRunePaths[ri.FullPath] = true
			toExpand = append(toExpand, ri)
		}
		if len(toExpand) == 0 {
			break
		}

		fmt.Printf("\n📤 Dispatching %d expansions...\n", len(toExpand))

		results := make([]expansionResult, len(toExpand))
		var wg sync.WaitGroup
		var totalReqNanos int64
		levelStart := time.Now()

		for i, ri := range toExpand {
			wg.Add(1)
			go func(i int, ri RuneExpansionInfo) {
				defer wg.Done()

				extendedReq := fmt.Sprintf(`Forget the prior decomposition. Imagine you are seeing "%s" for the first time, in isolation, as a black-box function you have to implement.

Question: what 0–3 PRIVATE helper functions would you write inside "%s"'s body to do its job? Helpers that no other function would ever call. Implementation details only.

Call the decompose tool. The runes map keys must be of the form "%s.<new_helper_name>". Example, for a different rune: if you were expanding "image.compress", reasonable helpers would be "image.compress.detect_format", "image.compress.choose_quality", "image.compress.encode_bytes". Each is a verb-phrase describing one internal step.

If "%s" is a single primitive operation (like an arithmetic op or a single syscall) and would have no private helpers in its body, return an empty runes map ({}). That is the correct answer.

Hard rules:
- Reply ONLY by calling the decompose tool.
- Never include sibling-level functions, never repeat existing names, never include "%s" itself.
- At most 3 helpers.`, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath)

				localMsgs := make([]openai.ChatCompletionMessage, 0, len(baseMessages)+1)
				localMsgs = append(localMsgs, baseMessages...)
				localMsgs = append(localMsgs, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: extendedReq,
				})

				reqStart := time.Now()
				response, err := client.Decompose(ctx, &localMsgs)
				reqDur := time.Since(reqStart)
				atomic.AddInt64(&totalReqNanos, int64(reqDur))
				dur := reqDur.Round(time.Millisecond)

				if err != nil {
					stdoutMu.Lock()
					fmt.Printf("   ⚠️  %s: %v (%s)\n", ri.FullPath, err, dur)
					stdoutMu.Unlock()
					results[i] = expansionResult{ri, nil, err}
					return
				}

				respVal, ok := response.(DecompositionResponse)
				if !ok {
					stdoutMu.Lock()
					if clar, isClar := response.(ClarificationRequest); isClar {
						fmt.Printf("   ⚠️  %s: model returned text instead of tool call (%s): %q\n", ri.FullPath, dur, clar.Message)
					} else {
						fmt.Printf("   ⚠️  %s: unexpected response type %T (%s): %+v\n", ri.FullPath, response, dur, response)
					}
					stdoutMu.Unlock()
					results[i] = expansionResult{ri, nil, fmt.Errorf("unexpected response type %T", response)}
					return
				}
				if respVal.ProjectPackage.Name == "" {
					stdoutMu.Lock()
					fmt.Printf("   ⚠️  %s: tool call had empty project_package.name (%s)\n      parsed response: %+v\n", ri.FullPath, dur, respVal)
					stdoutMu.Unlock()
					results[i] = expansionResult{ri, nil, fmt.Errorf("empty project_package.name")}
					return
				}

				newRunes := collectRunesForExpansion(&respVal)
				stdoutMu.Lock()
				if len(newRunes) == 0 {
					fmt.Printf("   ✓ %s: leaf (%s)\n", ri.FullPath, dur)
				} else {
					fmt.Printf("   ➜ %s: %d sub-runes (%s)\n", ri.FullPath, len(newRunes), dur)
				}
				stdoutMu.Unlock()

				results[i] = expansionResult{ri, &respVal, nil}
			}(i, ri)
		}

		wg.Wait()

		levelDur := time.Since(levelStart)
		sumDur := time.Duration(atomic.LoadInt64(&totalReqNanos))
		factor := float64(sumDur) / float64(levelDur)
		fmt.Printf("   ⏱️  level wall-clock: %s, sum of %d requests: %s (parallelism factor: %.1fx)\n",
			levelDur.Round(time.Millisecond),
			len(toExpand),
			sumDur.Round(time.Millisecond),
			factor,
		)

		var nextLevel []RuneExpansionInfo
		for _, r := range results {
			if r.resp == nil {
				continue
			}

			newRunes := collectRunesForExpansion(r.resp)
			if len(newRunes) == 0 {
				continue
			}

			childDecomposition := &AutoDecomposition{
				Path:       r.runeInfo.FullPath,
				Depth:      r.runeInfo.Depth + 1,
				Response:   r.resp,
				ParentPath: "",
				ChildPaths: make([]string, 0),
			}
			allDecompositions = append(allDecompositions, childDecomposition)

			if r.runeInfo.ParentDecomposition != nil {
				r.runeInfo.ParentDecomposition.ChildPaths = append(r.runeInfo.ParentDecomposition.ChildPaths, r.runeInfo.FullPath)
				childDecomposition.ParentPath = r.runeInfo.ParentDecomposition.Path
			}

			for j := range newRunes {
				newRunes[j].Depth = r.runeInfo.Depth + 1
				newRunes[j].ParentDecomposition = childDecomposition
			}
			nextLevel = append(nextLevel, newRunes...)
			totalRunesCount += countTotalRunes(r.resp)
		}

		currentLevel = nextLevel
	}

	fmt.Printf("\n")
	printCompleteTree(allDecompositions, "root", 0, true)

	separator := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered (max depth %d)\n", len(allDecompositions), totalRunesCount, cfg.MaxDepth)
	fmt.Printf("%s\n", separator)
}

func (c *Client) Decompose(ctx context.Context, messages *[]openai.ChatCompletionMessage) (any, error) {
	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:      MODEL_NAME,
		Messages:   *messages,
		Tools:      []openai.Tool{decomposeTool},
		ToolChoice: "auto",
	})
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	msg := resp.Choices[0].Message
	*messages = append(*messages, msg)

	if len(msg.ToolCalls) > 0 {
		call := msg.ToolCalls[0]
		if call.Function.Name != "decompose" {
			return nil, fmt.Errorf("unexpected tool call: %s", call.Function.Name)
		}
		var decomp DecompositionResponse
		if err := json.Unmarshal([]byte(call.Function.Arguments), &decomp); err != nil {
			return nil, fmt.Errorf("parsing decompose arguments: %w (raw: %s)", err, call.Function.Arguments)
		}
		*messages = append(*messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			ToolCallID: call.ID,
			Content:    "decomposition recorded",
		})
		return decomp, nil
	}

	if strings.TrimSpace(msg.Content) != "" {
		return ClarificationRequest{Message: msg.Content}, nil
	}

	return nil, fmt.Errorf("model returned neither a tool call nor content")
}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	if resp == nil || resp.ProjectPackage.Name == "" {
		return runes
	}

	if len(resp.ProjectPackage.Runes) > 0 {
		for name := range resp.ProjectPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.ProjectPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.ProjectPackage.Name, name)
			}
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		for name := range resp.StdPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.StdPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.StdPackage.Name, name)
			}
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	return runes
}

func countTotalRunes(resp *DecompositionResponse) int {
	if resp == nil {
		return 0
	}
	count := 0
	if len(resp.ProjectPackage.Runes) > 0 {
		count += len(resp.ProjectPackage.Runes)
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		count += len(resp.StdPackage.Runes)
	}
	return count
}

func printCompleteTree(allDecompositions []*AutoDecomposition, path string, depth int, isRoot bool) {
	var decomposition *AutoDecomposition
	for _, d := range allDecompositions {
		if d.Path == path {
			decomposition = d
			break
		}
	}

	if decomposition == nil || decomposition.Response == nil {
		return
	}

	resp := decomposition.Response

	if isRoot {
		fmt.Printf("🌳 ROOT DECOMPOSITION: %s\n", path)
	} else {
		indent := strings.Repeat("   ", depth)
		fmt.Printf("%s🔸 EXPANDED: %s\n", indent, path)
	}

	if len(resp.ProjectPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📦 %s", resp.ProjectPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.ProjectPackage.Runes, depth+1)
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📚 %s", resp.StdPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.StdPackage.Runes, depth+1)
	}

	for _, childPath := range decomposition.ChildPaths {
		printCompleteTree(allDecompositions, childPath, depth+1, false)
	}
}

func printRunesIndented(runes map[string]Rune, indentLevel int) {
	if len(runes) == 0 {
		return
	}

	indent := strings.Repeat("   ", indentLevel)

	for name, rune := range runes {
		fmt.Printf("%s├─ %s\n", indent, name)
		if rune.Description != "" {
			descIndent := strings.Repeat("   ", indentLevel+1)
			wrappedDesc := wrapText(rune.Description, 70-len(descIndent))
			fmt.Printf("%s│  └─ %s\n", descIndent, wrappedDesc)
		}
		if rune.FunctionSig != "" {
			sigIndent := strings.Repeat("   ", indentLevel+1)
			fmt.Printf("%s│     sig: %s\n", sigIndent, rune.FunctionSig)
		}
	}
}

func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	return text[:maxWidth-3] + "..."
}
