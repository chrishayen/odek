package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client calls the Anthropic API directly.
type Client struct {
	Model   string
	Token   string
	BaseURL string
	Mock    bool
	http    *http.Client
}

// New creates a Client from config fields.
func New(model, token string, mock bool) *Client {
	if model == "" {
		model = "claude-sonnet-4-6"
	}
	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return &Client{
		Model:   model,
		Token:   token,
		BaseURL: baseURL,
		Mock:    mock,
		http: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
			},
		},
	}
}

// Call sends a system+user prompt to the Anthropic API and returns the text response.
func (c *Client) Call(systemPrompt, userPrompt string) (string, error) {
	if c.Mock {
		return mockResponse(systemPrompt, userPrompt), nil
	}

	body := map[string]any{
		"model":      c.Model,
		"max_tokens": 16384,
		"system":     systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.Token)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != 200 {
		snippet := string(respBody)
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return "", fmt.Errorf("api error %d: %s", resp.StatusCode, snippet)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return result.Content[0].Text, nil
}

// StripCodeFences removes markdown code fences from Claude output.
func StripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	for _, fence := range []string{"```typescript\n", "```ts\n", "```go\n", "```python\n", "```\n"} {
		s = strings.ReplaceAll(s, fence, "")
	}
	s = strings.ReplaceAll(s, "```", "")
	return strings.TrimSpace(s)
}

// mockResponse returns deterministic responses for e2e tests.
func mockResponse(systemPrompt, userPrompt string) string {
	if strings.Contains(systemPrompt, "decompose") || strings.Contains(systemPrompt, "composition tree") {
		return mockDecomposeResponse()
	}
	if strings.Contains(systemPrompt, "RESULT: ALL PASS") || strings.Contains(systemPrompt, "verify") {
		return "PASS + test case 1 — implementation matches\nRESULT: ALL PASS"
	}
	return mockHydrateResponse(userPrompt)
}

func mockDecomposeResponse() string {
	return `std
  std.auth
    @ (email: string, password: string) -> result[bool, string]
    + authenticates user with valid credentials
    - returns error for invalid credentials
    std.auth.validate_email
      @ (email: string) -> bool
      + Given 'user@example.com', returns true
      + Given 'a@b.co', returns true
      - Given an empty string, returns false
      - Given 'no-at-sign', returns false
      - Given '@missing-local.com', returns false
    std.auth.hash_password
      @ (password: string) -> result[string, string]
      + Given 'secret123', returns a non-empty string
      + Given the same input twice, returns the same hash both times
      - Given an empty string, returns an error

test_project
  @ () -> result[void, string]
  + user login flow completes successfully
  - returns error for invalid email
  test_project.login
    @ (email: string, password: string) -> result[bool, string]
    + authenticates and returns session
    - rejects invalid email before checking password
    -> std.auth.validate_email
    -> std.auth.hash_password
    test_project.login.store_credentials
      @ (user_id: string, hashed_password: string) -> result[bool, string]
      + Given valid user ID and hash, returns success
      - Given empty user ID, returns an error`
}

func mockHydrateResponse(prompt string) string {
	return fmt.Sprintf(`=== FILE: go.mod ===
module valkyrie-rune

go 1.22
=== END FILE ===

=== FILE: main.go ===
package main

import "fmt"

// Generated for: %s
func Run() string {
	return "Hello, World!"
}

func main() {
	fmt.Println(Run())
}
=== END FILE ===

=== FILE: main_test.go ===
package main

import "testing"

func TestRun(t *testing.T) {
	got := Run()
	want := "Hello, World!"
	if got != want {
		t.Errorf("Run() = %%q, want %%q", got, want)
	}
}
=== END FILE ===
`, prompt)
}
