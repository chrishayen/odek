package openai

import (
	"context"
	"fmt"
	"net/http"
)

// Client represents the client for interacting with an OpenAI-compatible API.
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewClient creates a new client for the API.
func NewClient(baseURL string, apiKey ...string) (*Client, error) {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:1234" // Default local dev URL
	}

	apiKeyVal := ""
	if len(apiKey) > 0 {
		apiKeyVal = apiKey[0]
	}

	return &Client{
		baseURL: fmt.Sprintf("%s/v1", baseURL),
		apiKey:  apiKeyVal,
		client:  &http.Client{},
	}, nil
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
