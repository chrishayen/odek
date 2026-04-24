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
		Model: openai.DefaultModel,
		Messages: []openai.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
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
	if decomposition.RuneTree == nil {
		return nil, fmt.Errorf("validation failed: rune_tree is required")
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
