package openai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents the client for interacting with an OpenAI-compatible API.
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewClient creates a new client for the API.
func NewClient(baseURL string, apiKey ...string) (*Client, error) {
	return NewClientWithHTTPClient(baseURL, &http.Client{Timeout: 10 * time.Minute}, apiKey...)
}

// NewClientWithHTTPClient creates a client using the provided HTTP client.
// It is useful for tests and callers that need custom transports, proxies, or
// timeout behavior.
func NewClientWithHTTPClient(baseURL string, httpClient *http.Client, apiKey ...string) (*Client, error) {
	normalizedBaseURL, err := normalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Minute}
	}

	apiKeyVal := ""
	if len(apiKey) > 0 {
		apiKeyVal = apiKey[0]
	}

	return &Client{
		baseURL: normalizedBaseURL,
		apiKey:  apiKeyVal,
		client:  httpClient,
	}, nil
}

func normalizeBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = "http://127.0.0.1:1234"
	}
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	raw = strings.TrimRight(raw, "/")

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid API base URL %q: %w", raw, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid API base URL %q: scheme and host are required", raw)
	}
	if !strings.HasSuffix(u.Path, "/v1") {
		u.Path = strings.TrimRight(u.Path, "/") + "/v1"
	}
	return u.String(), nil
}

// ModelInfo represents a model available on the server.
type ModelInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
	Created int64  `json:"created,omitempty"`
	Meta    struct {
		Dimension int `json:"dimension,omitempty"`
	} `json:"meta,omitempty"`
}

// ListModelsResponse represents the root of a list models response.
type ListModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ListModels fetches the list of available models from the API.
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) {
	var result ListModelsResponse
	if err := c.doJSON(ctx, http.MethodGet, "/models", nil, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetModelInfo fetches detailed information about a specific model.
func (c *Client) GetModelInfo(ctx context.Context, modelID string) (*ModelInfo, error) {
	var info ModelInfo
	if err := c.doJSON(ctx, http.MethodGet, "/models/"+modelID, nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// HealthCheck performs a simple GET request to /health to verify connectivity.
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.doJSON(ctx, http.MethodGet, "/health", nil, nil)
}
