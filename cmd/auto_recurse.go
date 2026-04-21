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
	BASE_URL      = "http://localhost:8080"
	EXAMPLES_DIR  = "examples"
	TOOL_LOG_PATH = "/tmp/odek-example-log.jsonl"
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

	api, err := openai.NewClient(BASE_URL)
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

	// Stage 1: estimate effort (retained for the effort label in the output;
	// the two-pass pipeline doesn't vary depth with effort anymore).
	fmt.Printf("\n🎯 Estimating effort for: %s\n", req)
	effortStart := time.Now()
	est, err := effortpkg.Estimate(ctx, api, req)
	if err != nil {
		fmt.Printf("⚠️  effort estimation failed: %v — defaulting to level 3\n", err)
		est = effortpkg.Result{Level: 3, Reason: "default (estimator failed)"}
	}
	cfg := decomposer.ConfigForEffort(est.Level)
	fmt.Printf("🎯 effort: %d/5 — %s (%s)\n", est.Level, est.Reason, time.Since(effortStart).Round(time.Millisecond))

	// Stage 2: run the two-pass decompose and stream events.
	fmt.Printf("\nDecomposing: %s...\n", req)
	runStart := time.Now()
	sess, err := dec.NewSession(ctx, req, est.Level, est.Reason, cfg, decomposer.SessionContext{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for evt := range sess.Events {
		printDecompositionEvent(evt)
	}

	resp := sess.Response()
	if resp == nil {
		fmt.Printf("\n⚠️  No decomposition produced: %s\n", sess.Snapshot().ErrorMsg)
		return
	}

	printCompleteTree(resp)

	total := countTreeRunes(resp)
	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Printf("📊 SUMMARY: %d runes (%s total)\n", total, time.Since(runStart).Round(time.Millisecond))
	fmt.Printf("%s\n", strings.Repeat("=", 70))
}

