package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	_ "embed"
	"shotgun.dev/odek/decompose"
	"shotgun.dev/odek/internal/tui"
	openai "shotgun.dev/odek/openai"
)

// System prompt embedded at compile time from decompose.md
//
//go:embed decompose.md
var systemPrompt string

func main() {
	// 1. Initialize the client
	// You can change this to point to a remote API (e.g., https://api.openai.com/v1)
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default local server
	}

	client, err := openai.NewClient(baseURL)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// 2. Check for -p flag and process prompt if present
	prompt := ""
	for i, arg := range os.Args {
		if (i > 0 && strings.HasPrefix(arg, "-p=")) || arg == "-p" {
			prompt = os.Args[i+1]
			break
		}
	}

	// 3. Check for -d flag
	var dValue string
	for i, arg := range os.Args {
		if (i > 0 && strings.HasPrefix(arg, "-d=")) || arg == "-d" {
			dValue = os.Args[i+1]
			break
		}
	}

	// 4. Check for -d flag (decompose feature)
	if dValue != "" {
		response, err := decompose.Decompose(ctx, client, systemPrompt, dValue)
		if err != nil {
			log.Fatalf("Decompose failed: %v", err)
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

		return
	}

	if prompt != "" {
		// Build chat request with the prompt
		request := &openai.ChatCompletionRequest{
			Model:    "default",
			Messages: []openai.ChatMessage{{Role: "user", Content: prompt}},
		}

		response, err := client.Chat(ctx, request)
		if err != nil {
			log.Fatalf("Chat failed: %v", err)
		}

		// Print the result
		for _, choice := range response.Choices {
			fmt.Printf("\n=== Response ===\n%s\n", choice.Message.Content)
			if response.Usage != nil {
				fmt.Printf("Tokens: prompt=%d, completion=%d, total=%d\n",
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens,
					response.Usage.TotalTokens)
			}
		}

		return // Exit after processing prompt
	}

	// 3. Launch TUI when no arguments provided
	tui.Run()
}
