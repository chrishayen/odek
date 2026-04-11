package decompose

import (
	"context"
	"fmt"

	openai "shotgun.dev/odek/openai"
)

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
