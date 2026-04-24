package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"shotgun.dev/odek/internal/decomposer"
	effortpkg "shotgun.dev/odek/internal/effort"
	openai "shotgun.dev/odek/openai"
)

const (
	DEFAULT_BASE_URL = "http://localhost:8080"
	EXAMPLES_DIR     = "examples"
	TOOL_LOG_PATH    = "/tmp/odek-example-log.jsonl"
)

func getRequirement(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)
	return requirement, nil
}

func main() {
	ctx := context.Background()
	printBanner()

	api, err := openai.NewClient(apiBaseURL(), apiKey())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create openai client: %v\n", err)
		os.Exit(1)
	}
	dec, err := decomposer.NewDecomposer(api, EXAMPLES_DIR, TOOL_LOG_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create decomposer: %v\n", err)
		os.Exit(1)
	}

	stdin := bufio.NewReader(os.Stdin)
	req, err := getRequirement(stdin)
	if err != nil || req == "" {
		fmt.Println("No requirement provided. Exiting.")
		os.Exit(1)
	}

	// Stage 1: estimate effort
	fmt.Printf("\n🎯 Estimating effort for: %s\n", req)
	effortStart := time.Now()
	est, err := effortpkg.Estimate(ctx, api, req)
	if err != nil {
		fmt.Printf("⚠️  effort estimation failed: %v — defaulting to level 3\n", err)
		est = effortpkg.Result{Level: 3, Reason: "default (estimator failed)"}
	}
	cfg := decomposer.ConfigForEffort(est.Level)
	fmt.Printf("🎯 effort: %d/5 — %s (%s)\n", est.Level, est.Reason, time.Since(effortStart).Round(time.Millisecond))
	fmt.Printf("   config: parallel=%d depth=%d cap=%d recurse=%v\n",
		cfg.ParallelInitial, cfg.MaxDepth, cfg.RuneCap, cfg.Recurse)

	// Stages 2 & 3: initial decomposition (single or N-way + merge)
	fmt.Printf("\nDecomposing: %s...\n", req)
	initStart := time.Now()
	sess, err := dec.NewSession(ctx, req, est.Level, est.Reason, cfg, decomposer.SessionContext{})
	fmt.Printf("⏱️  initial decompose: %s\n", time.Since(initStart).Round(time.Millisecond))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// Stage 4: present + confirm
	printInitialDecomposition(sess.Root.Response)

	if !cfg.Recurse {
		fmt.Println("\n(Skipping recursion: requirement is trivial enough.)")
		return
	}
	if len(sess.TopLevelPaths()) == 0 {
		fmt.Println("\n(No runes to expand.)")
		return
	}

	if !confirm(stdin, fmt.Sprintf("\nProceed with recursion (max depth %d, max %d runes)? [y/N] ", cfg.MaxDepth, cfg.RuneCap)) {
		fmt.Println("Stopping before recursion. Initial decomposition is above.")
		return
	}

	// Stage 5: recurse
	fmt.Printf("\n🔄 Starting auto-recursion\n")
	for evt := range dec.ExpandStreaming(ctx, sess, cfg) {
		printExpansionEvent(evt)
	}
	printCompleteTree(sess.AllDecompositions(), "root", 0, true)

	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered (max depth %d)\n",
		len(sess.AllDecompositions()), sess.Snapshot().TotalRunes, cfg.MaxDepth)
	fmt.Printf("%s\n", strings.Repeat("=", 70))
}

func apiBaseURL() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return DEFAULT_BASE_URL
}

func apiKey() string {
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		return v
	}
	return os.Getenv("API_KEY")
}

func confirm(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}
