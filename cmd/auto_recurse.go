package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	effortpkg "shotgun.dev/odek/internal/effort"
	"shotgun.dev/odek/openai"
)

const (
	BASE_URL            = "http://localhost:8080"
	EXAMPLES_DIR        = "examples"
	TOOL_LOG_PATH       = "/tmp/odek-example-log.jsonl"
	MAX_TOOL_ITERATIONS = 6
)

//go:embed decompose.md
var SYSTEM_PROMPT string

type RunConfig struct {
	ParallelInitial int
	MaxDepth        int
	RuneCap         int
	Recurse         bool
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
	est, err := effortpkg.Estimate(ctx, client.api, req)
	if err != nil {
		fmt.Printf("⚠️  effort estimation failed: %v — defaulting to level 3\n", err)
		est = effortpkg.Result{Level: 3, Reason: "default (estimator failed)"}
	}
	cfg := configForEffort(est.Level)
	fmt.Printf("🎯 effort: %d/5 — %s (%s)\n", est.Level, est.Reason, time.Since(effortStart).Round(time.Millisecond))
	fmt.Printf("   config: parallel=%d depth=%d cap=%d recurse=%v\n",
		cfg.ParallelInitial, cfg.MaxDepth, cfg.RuneCap, cfg.Recurse)

	// Stages 2 & 3: initial decompose (single or N-way + merge)
	response, baseMessages, err := runInitialDecompose(ctx, req, cfg)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	var rootResponse DecompositionResponse
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

// runInitialDecompose produces the root decomposition — either a single pass
// or N parallel attempts merged into a consensus. Returns the response (either
// a DecompositionResponse or a ClarificationRequest) and the conversation
// history the caller should carry into the recursion phase.
func runInitialDecompose(ctx context.Context, req string, cfg RunConfig) (any, []openai.ChatMessage, error) {
	if cfg.ParallelInitial == 1 {
		baseMessages := newConversation(req)
		fmt.Printf("\nDecomposing: %s...\n", req)
		initStart := time.Now()
		response, history, err := client.Decompose(ctx, baseMessages)
		fmt.Printf("⏱️  initial decompose: %s\n", time.Since(initStart).Round(time.Millisecond))
		return response, history, err
	}

	fmt.Printf("\n🚀 Running %d parallel initial decompositions...\n", cfg.ParallelInitial)
	parallelStart := time.Now()
	attempts := parallelInitialDecompose(ctx, req, cfg.ParallelInitial)
	fmt.Printf("✅ %d/%d attempts succeeded (%s)\n",
		len(attempts), cfg.ParallelInitial, time.Since(parallelStart).Round(time.Millisecond))

	if len(attempts) == 0 {
		return nil, nil, fmt.Errorf("all parallel attempts failed")
	}
	if len(attempts) == 1 {
		return attempts[0], newConversation(req), nil
	}

	fmt.Printf("\n🔀 Merging %d attempts...\n", len(attempts))
	mergeStart := time.Now()
	merged, mergedMsgs, err := client.MergeAttempts(ctx, req, attempts)
	fmt.Printf("⏱️  merge: %s\n", time.Since(mergeStart).Round(time.Millisecond))
	if err != nil {
		fmt.Printf("⚠️  merge failed: %v — using first attempt\n", err)
		return attempts[0], newConversation(req), nil
	}
	return merged, mergedMsgs, nil
}

func configForEffort(level int) RunConfig {
	switch level {
	case 1:
		return RunConfig{ParallelInitial: 1, MaxDepth: 0, RuneCap: 10, Recurse: false}
	case 2:
		return RunConfig{ParallelInitial: 1, MaxDepth: 10, RuneCap: 25, Recurse: true}
	case 3:
		return RunConfig{ParallelInitial: 3, MaxDepth: 10, RuneCap: 50, Recurse: true}
	case 4:
		return RunConfig{ParallelInitial: 5, MaxDepth: 10, RuneCap: 100, Recurse: true}
	case 5:
		return RunConfig{ParallelInitial: 5, MaxDepth: 10, RuneCap: 200, Recurse: true}
	}
	return RunConfig{ParallelInitial: 3, MaxDepth: 10, RuneCap: 50, Recurse: true}
}

func confirm(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}
