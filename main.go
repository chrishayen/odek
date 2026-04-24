package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"shotgun.dev/odek/internal/decomposer"
	"shotgun.dev/odek/internal/tui"
	openai "shotgun.dev/odek/openai"
)

const (
	examplesDir    = "examples"
	toolLogPath    = "/tmp/odek-example-log.jsonl"
	defaultBaseURL = "http://localhost:8080"
)

type cliOptions struct {
	prompt       string
	decomposeReq string
	jsonOutput   bool
}

func main() {
	opts, err := parseFlags(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if opts.prompt != "" && opts.decomposeReq != "" {
		fmt.Fprintln(os.Stderr, "use either -p for chat or -d for decomposition, not both")
		os.Exit(2)
	}

	client, err := newAPIClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	if opts.decomposeReq != "" {
		dec, err := decomposer.NewDecomposer(client, examplesDir, toolLogPath)
		if err != nil {
			log.Fatalf("Failed to create decomposer: %v", err)
		}
		if err := runDirectDecompose(ctx, dec, opts.decomposeReq); err != nil {
			log.Fatalf("Decompose failed: %v", err)
		}
		return
	}

	if opts.prompt != "" {
		if err := runDirectPrompt(ctx, client, opts.prompt, opts.jsonOutput); err != nil {
			log.Fatalf("Chat failed: %v", err)
		}
		return
	}

	dec, err := decomposer.NewDecomposer(client, examplesDir, toolLogPath)
	if err != nil {
		log.Printf("decomposer init failed: %v — /decompose will be unavailable", err)
	}
	if err := tui.Run(ctx, client, dec); err != nil {
		log.Fatalf("TUI failed: %v", err)
	}
}

func parseFlags(args []string) (cliOptions, error) {
	var opts cliOptions
	fs := flag.NewFlagSet("odek", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&opts.prompt, "p", "", "send one prompt and print the model response")
	fs.StringVar(&opts.decomposeReq, "d", "", "decompose one requirement and print the structured rune tree as JSON")
	fs.BoolVar(&opts.jsonOutput, "j", false, "print direct chat responses as raw JSON")
	fs.BoolVar(&opts.jsonOutput, "json", false, "print direct chat responses as raw JSON")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if fs.NArg() > 0 {
		return opts, fmt.Errorf("unexpected positional arguments: %v", fs.Args())
	}
	return opts, nil
}

func newAPIClient() (*openai.Client, error) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}
	return openai.NewClient(baseURL, apiKey)
}

func runDirectDecompose(ctx context.Context, dec *decomposer.Decomposer, requirement string) error {
	cfg := decomposer.ConfigForEffort(2)
	cfg.ParallelInitial = 1
	cfg.MaxDepth = 0
	cfg.RuneCap = 0
	cfg.Recurse = false

	sess, err := dec.NewSession(ctx, requirement, 2, "direct CLI decomposition", cfg, decomposer.SessionContext{})
	if err != nil {
		return err
	}
	if sess == nil || sess.Root == nil || sess.Root.Response == nil {
		return fmt.Errorf("decomposer returned no root response")
	}
	jsonOutput, err := json.MarshalIndent(sess.Root.Response, "", "  ")
	if err != nil {
		return fmt.Errorf("format response JSON: %w", err)
	}
	fmt.Println(string(jsonOutput))
	return nil
}

func runDirectPrompt(ctx context.Context, client *openai.Client, prompt string, jsonOutput bool) error {
	request := &openai.ChatCompletionRequest{
		Model:    openai.DefaultModel,
		Messages: []openai.ChatMessage{{Role: openai.RoleUser, Content: prompt}},
	}

	response, err := client.Chat(ctx, request)
	if err != nil {
		return err
	}

	if jsonOutput {
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("format response JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	for _, choice := range response.Choices {
		fmt.Printf("\n=== Response ===\n%s\n", choice.Message.Content)
		if response.Usage != nil {
			fmt.Printf("Tokens: prompt=%d, completion=%d, total=%d\n",
				response.Usage.PromptTokens,
				response.Usage.CompletionTokens,
				response.Usage.TotalTokens)
		}
	}
	return nil
}
