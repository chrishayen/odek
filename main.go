package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "shotgun.dev/odek/openai"
)

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

	// 3. List available models (original behavior)
	fmt.Println("Fetching available models...")
	modelsList, err := client.ListModels(ctx)
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	// 4. Print results
	fmt.Printf("Found %d models:\n\n", len(modelsList))
	for _, m := range modelsList {
		fmt.Printf("- ID: %s | Name: %s\n", m.ID, m.Name)
	}

	// Optional: Check health
	if err := client.HealthCheck(ctx); err != nil {
		log.Printf("Health check warning (might be expected if API doesn't have /health): %v", err)
	}
}
