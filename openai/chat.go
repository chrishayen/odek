package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Event represents a streaming event in the API response.
type Event struct {
	Type string `json:"type"`
	Data struct {
		Message  string      `json:"message,omitempty"`
		Choice   Choice      `json:"choice,omitempty"`
		MetaData interface{} `json:"meta_data,omitempty"`
	} `json:"data"`
}

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// Tool describes a function the model is allowed to call.
type Tool struct {
	Type     string              `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

// FunctionDefinition is the JSON-schema description of a callable function.
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

// ToolCall is a single function invocation emitted by the model.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction carries the function name plus its JSON-encoded arguments.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatCompletionRequest represents the request body for chat completions.
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Tools       []Tool        `json:"tools,omitempty"`
	ToolChoice  any           `json:"tool_choice,omitempty"`
}

// Choice represents a single completion option.
type Choice struct {
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Score        float64     `json:"score,omitempty"`
	Index        int         `json:"index"`
	Delta        Delta       `json:"delta,omitempty"`
	Details      *Usage      `json:"details,omitempty"`
}

// Delta represents the delta for a choice in streaming responses.
type Delta struct {
	Content string `json:"content,omitempty"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse represents the API response for chat completions.
type ChatCompletionResponse struct {
	ID                string   `json:"id,omitempty"`
	Object            string   `json:"object,omitempty"`
	Created           int64    `json:"created,omitempty"`
	Model             string   `json:"model,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// SystemPromptBuilder constructs the system message for conversation setup.
type SystemPromptBuilder struct {
	prompt   string
	metadata map[string]string
}

// NewSystemPromptBuilder creates a new builder with an optional base prompt.
func NewSystemPromptBuilder(base ...string) *SystemPromptBuilder {
	return &SystemPromptBuilder{
		prompt:   "",
		metadata: make(map[string]string),
	}
}

// SetBase sets the primary system instruction text.
func (sb *SystemPromptBuilder) SetBase(text string) *SystemPromptBuilder {
	sb.prompt = text
	return sb
}

// AddMetadata adds key-value pairs for model context.
func (sb *SystemPromptBuilder) AddMetadata(key, value string) *SystemPromptBuilder {
	sb.metadata[key] = value
	return sb
}

// Build returns the complete system prompt message as a ChatMessage.
func (sb *SystemPromptBuilder) Build() ChatMessage {
	metadata := ""
	if len(sb.metadata) > 0 {
		var parts []string
		for k, v := range sb.metadata {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		metadata = " | " + joinParts(parts, ", ")
	}
	return ChatMessage{
		Role:    "system",
		Content: sb.prompt + metadata,
		Name:    "system",
	}
}

// Chat completes a conversation with the AI model. Supports both single-turn and multi-turn chat.
func (c *Client) Chat(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := c.baseURL + "/chat/completions"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func joinParts(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}
