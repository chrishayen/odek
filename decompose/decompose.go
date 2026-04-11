package decompose

import (
	"context"
	"fmt"

	openai "shotgun.dev/odek/openai"
)

// DecomposeResult contains both the raw API response and parsed structured data
type DecomposeResult struct {
	Response      *openai.ChatCompletionResponse
	Decomposition *Decomposition
}

// Decompose sends a chat request with system prompt + user message to the API.
func Decompose(ctx context.Context, client *openai.Client, systemPrompt, userMessage string) (*openai.ChatCompletionResponse, error) {
	request := &openai.ChatCompletionRequest{
		Model: "default",
		Messages: []openai.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
	}

	fmt.Printf("DEBUG decompose: Building request with %d messages\n", len(request.Messages))
	for i, msg := range request.Messages {
		fmt.Printf("DEBUG decompose: Message %d (role=%s): %.80s...\n", i, msg.Role, msg.Content)
	}

	return client.Chat(ctx, request)
}

// DecomposeStructured sends a chat request and returns both raw response and parsed structured data.
func DecomposeStructured(ctx context.Context, client *openai.Client, systemPrompt, userMessage string) (*DecomposeResult, error) {
	response, err := Decompose(ctx, client, systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("decompose failed: %w", err)
	}

	var decomposition *Decomposition
	for _, choice := range response.Choices {
		dec, parseErr := ParseDecomposition(choice.Message.Content)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse structured output: %w", parseErr)
		}
		decomposition = dec
		break
	}

	if decomposition == nil {
		return nil, fmt.Errorf("no response content to parse")
	}

	// Validate the decomposition structure
	if err := decomposition.RuneTree.Validate("root"); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &DecomposeResult{
		Response:      response,
		Decomposition: decomposition,
	}, nil
}
