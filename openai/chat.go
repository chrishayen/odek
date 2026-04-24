package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
)

// thinkingCallbackKey is the ctx key for a per-request reasoning_content
// callback. When set, Chat transparently switches to SSE streaming so the
// callback can fire on every reasoning delta the server emits.
type thinkingCallbackKey struct{}

// WithThinkingCallback attaches a reasoning_content delta callback to ctx.
// Any Chat call made with this context streams, and cb is invoked once per
// reasoning_content chunk the server sends. Pass nil to clear.
func WithThinkingCallback(ctx context.Context, cb func(string)) context.Context {
	if cb == nil {
		return ctx
	}
	return context.WithValue(ctx, thinkingCallbackKey{}, cb)
}

func thinkingCallbackFromContext(ctx context.Context) func(string) {
	v, _ := ctx.Value(thinkingCallbackKey{}).(func(string))
	return v
}

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
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	Name             string     `json:"name,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty"`
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
	Stream      bool          `json:"stream,omitempty"`
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
	Content          string `json:"content,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
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
	prompt := ""
	if len(base) > 0 {
		prompt = base[0]
	}
	return &SystemPromptBuilder{
		prompt:   prompt,
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
		sort.Strings(parts)
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
//
// If the ctx carries a thinking callback (via WithThinkingCallback), Chat
// switches to SSE streaming so reasoning_content deltas can be surfaced
// live. Otherwise it uses the non-streaming endpoint. The final accumulated
// response is returned either way, so callers don't care which path ran.
func (c *Client) Chat(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if cb := thinkingCallbackFromContext(ctx); cb != nil || request.Stream {
		return c.chatStream(ctx, request, cb)
	}
	var response ChatCompletionResponse
	if err := c.doJSON(ctx, http.MethodPost, "/chat/completions", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// chatStream is the SSE path. It posts with stream=true, reassembles deltas
// into a single ChatCompletionResponse (so the caller sees the same shape),
// and fires thinkingCb on every reasoning_content chunk the server emits.
// thinkingCb may be nil — the stream still runs (useful when the caller set
// Stream=true explicitly).
func (c *Client) chatStream(ctx context.Context, request *ChatCompletionRequest, thinkingCb func(string)) (*ChatCompletionResponse, error) {
	reqCopy := *request
	reqCopy.Stream = true

	bodyBytes, err := json.Marshal(&reqCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	final := &ChatCompletionResponse{
		Choices: []Choice{{Message: ChatMessage{Role: RoleAssistant}}},
	}
	acc := &final.Choices[0]
	toolsByIdx := make(map[int]*ToolCall)
	var toolOrder []int

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}
		payload := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		if len(payload) == 0 || bytes.Equal(payload, []byte("[DONE]")) {
			if bytes.Equal(payload, []byte("[DONE]")) {
				break
			}
			continue
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Role             string `json:"role"`
					Content          string `json:"content"`
					ReasoningContent string `json:"reasoning_content"`
					ToolCalls        []struct {
						Index    int    `json:"index"`
						ID       string `json:"id"`
						Type     string `json:"type"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
			Usage *Usage `json:"usage"`
		}
		if err := json.Unmarshal(payload, &chunk); err != nil {
			continue
		}
		if chunk.Usage != nil {
			final.Usage = chunk.Usage
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		ch0 := chunk.Choices[0]
		if ch0.Delta.Content != "" {
			acc.Message.Content += ch0.Delta.Content
		}
		if ch0.Delta.ReasoningContent != "" {
			acc.Message.ReasoningContent += ch0.Delta.ReasoningContent
			if thinkingCb != nil {
				thinkingCb(ch0.Delta.ReasoningContent)
			}
		}
		for _, tc := range ch0.Delta.ToolCalls {
			existing, ok := toolsByIdx[tc.Index]
			if !ok {
				existing = &ToolCall{}
				toolsByIdx[tc.Index] = existing
				toolOrder = append(toolOrder, tc.Index)
			}
			if tc.ID != "" {
				existing.ID = tc.ID
			}
			if tc.Type != "" {
				existing.Type = tc.Type
			}
			if tc.Function.Name != "" {
				existing.Function.Name += tc.Function.Name
			}
			if tc.Function.Arguments != "" {
				existing.Function.Arguments += tc.Function.Arguments
			}
		}
		if ch0.FinishReason != "" {
			acc.FinishReason = ch0.FinishReason
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stream read: %w", err)
	}

	sort.Ints(toolOrder)
	for _, idx := range toolOrder {
		tc := toolsByIdx[idx]
		if tc.Type == "" {
			tc.Type = ToolTypeFunction
		}
		acc.Message.ToolCalls = append(acc.Message.ToolCalls, *tc)
	}

	return final, nil
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
// Each iteration posts the current history to /chat/completions with the given
// toolChoice (pass nil for the default "auto"). If the assistant turn has no tool
// calls, the loop returns that message as final. If it has tool calls, each one
// is dispatched to handler in order. A handler that returns terminal=true
// short-circuits the remaining calls in that turn; the loop then returns the
// assistant turn as final. Exceeding maxIterations without a terminal tool call
// or text-only response returns an error.
//
// The returned history always contains the full accumulated message list, including
// on error, so callers can log or feed it forward regardless of outcome.
func (c *Client) AskToolLoop(
	ctx context.Context,
	messages []ChatMessage,
	tools []Tool,
	handler ToolHandler,
	maxIterations int,
	toolChoice any,
) (final ChatMessage, history []ChatMessage, err error) {
	history = make([]ChatMessage, 0, len(messages)+maxIterations*2)
	history = append(history, messages...)

	tc := toolChoice
	if tc == nil {
		tc = "auto"
	}

	for range maxIterations {
		resp, chatErr := c.Chat(ctx, &ChatCompletionRequest{
			Model:      DefaultModel,
			Messages:   history,
			Tools:      tools,
			ToolChoice: tc,
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
