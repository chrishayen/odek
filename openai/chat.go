package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultModel = "default"

	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"

	ToolTypeFunction = "function"
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
		Role:    RoleSystem,
		Content: sb.prompt + metadata,
		Name:    "system",
	}
}

// doJSON sends method+path with an optional JSON body and decodes the response into out.
// Pass body=nil to omit the request body; pass out=nil to discard the response body.
// Treats 200 and 202 as success; everything else returns an error with the body attached.
func (c *Client) doJSON(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

// Chat completes a conversation with the AI model. Supports both single-turn and multi-turn chat.
func (c *Client) Chat(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var response ChatCompletionResponse
	if err := c.doJSON(ctx, http.MethodPost, "/chat/completions", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Ask sends a system+user prompt and returns the assistant's text content.
// An empty systemPrompt is omitted from the conversation.
func (c *Client) Ask(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	return c.AskMessages(ctx, buildPromptMessages(systemPrompt, userMessage))
}

// AskMessages sends a full conversation and returns the assistant's text content.
// Returns an error if the response has no choices.
func (c *Client) AskMessages(ctx context.Context, messages []ChatMessage) (string, error) {
	resp, err := c.Chat(ctx, &ChatCompletionRequest{
		Model:    DefaultModel,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	return resp.Choices[0].Message.Content, nil
}

// AskTool forces a single tool call and returns the raw JSON arguments string.
// Caller unmarshals into its own result type. Returns an error if no tool call was produced.
func (c *Client) AskTool(ctx context.Context, systemPrompt, userMessage string, tool Tool) (string, error) {
	if tool.Function == nil {
		return "", errors.New("tool.Function must not be nil")
	}
	resp, err := c.Chat(ctx, &ChatCompletionRequest{
		Model:    DefaultModel,
		Messages: buildPromptMessages(systemPrompt, userMessage),
		Tools:    []Tool{tool},
		ToolChoice: map[string]any{
			"type": ToolTypeFunction,
			"function": map[string]any{
				"name": tool.Function.Name,
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	calls := resp.Choices[0].Message.ToolCalls
	if len(calls) == 0 {
		return "", fmt.Errorf("no tool call in response (forced tool was %q)", tool.Function.Name)
	}
	return calls[0].Function.Arguments, nil
}

// ToolHandler is called once per tool call within an assistant turn during AskToolLoop.
// Return (result, false, nil) to append a tool-result message and continue the loop.
// Return (result, true, nil) to append the tool-result message and exit with this turn as terminal.
// Return (_, _, err) to abort the loop with the error; no tool result is appended.
type ToolHandler func(ctx context.Context, call ToolCall) (result string, terminal bool, err error)

// AskToolLoop runs a multi-turn tool-calling loop up to maxIterations.
//
// Each iteration posts the current history to /chat/completions with ToolChoice="auto".
// If the assistant turn has no tool calls, the loop returns that message as final.
// If it has tool calls, each one is dispatched to handler in order. A handler that
// returns terminal=true short-circuits the remaining calls in that turn; the loop
// then returns the assistant turn as final. Exceeding maxIterations without a
// terminal tool call or text-only response returns an error.
//
// The returned history always contains the full accumulated message list, including
// on error, so callers can log or feed it forward regardless of outcome.
func (c *Client) AskToolLoop(
	ctx context.Context,
	messages []ChatMessage,
	tools []Tool,
	handler ToolHandler,
	maxIterations int,
) (final ChatMessage, history []ChatMessage, err error) {
	history = make([]ChatMessage, 0, len(messages)+maxIterations*2)
	history = append(history, messages...)

	for range maxIterations {
		resp, chatErr := c.Chat(ctx, &ChatCompletionRequest{
			Model:      DefaultModel,
			Messages:   history,
			Tools:      tools,
			ToolChoice: "auto",
		})
		if chatErr != nil {
			return ChatMessage{}, history, fmt.Errorf("chat completion failed: %w", chatErr)
		}
		if len(resp.Choices) == 0 {
			return ChatMessage{}, history, errors.New("no choices in response")
		}

		msg := resp.Choices[0].Message
		history = append(history, msg)

		if len(msg.ToolCalls) == 0 {
			return msg, history, nil
		}

		terminated := false
		for _, call := range msg.ToolCalls {
			result, terminal, herr := handler(ctx, call)
			if herr != nil {
				return msg, history, herr
			}
			history = append(history, ChatMessage{
				Role:       RoleTool,
				ToolCallID: call.ID,
				Content:    result,
			})
			if terminal {
				terminated = true
				break
			}
		}
		if terminated {
			return msg, history, nil
		}
	}
	return ChatMessage{}, history, fmt.Errorf("exceeded %d tool iterations without terminal", maxIterations)
}

func buildPromptMessages(systemPrompt, userMessage string) []ChatMessage {
	msgs := make([]ChatMessage, 0, 2)
	if systemPrompt != "" {
		msgs = append(msgs, ChatMessage{Role: RoleSystem, Content: systemPrompt})
	}
	msgs = append(msgs, ChatMessage{Role: RoleUser, Content: userMessage})
	return msgs
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
