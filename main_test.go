package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	openai "shotgun.dev/odek/openai"
)

func TestRunDirectPromptWithOpenAICompatibleTransport(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var body strings.Builder
		err := json.NewEncoder(&body).Encode(openai.ChatCompletionResponse{
			Choices: []openai.Choice{{
				Index: 0,
				Message: openai.ChatMessage{
					Role:    openai.RoleAssistant,
					Content: "stubbed ok",
				},
				FinishReason: "stop",
			}},
			Usage: &openai.Usage{
				PromptTokens:     1,
				CompletionTokens: 2,
				TotalTokens:      3,
			},
		})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body.String())),
		}, nil
	})

	client, err := openai.NewClientWithHTTPClient("http://api.test", &http.Client{Transport: transport})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	output := captureStdout(t, func() {
		if err := runDirectPrompt(context.Background(), client, "hello", false); err != nil {
			t.Fatalf("runDirectPrompt: %v", err)
		}
	})

	if !strings.Contains(output, "stubbed ok") {
		t.Fatalf("output = %q, want assistant content", output)
	}
	if !strings.Contains(output, "Tokens: prompt=1, completion=2, total=3") {
		t.Fatalf("output = %q, want token usage", output)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return string(data)
}
