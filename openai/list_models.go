package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		baseURL = "http://localhost:8080" // Default local dev URL
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
	url := c.baseURL + "/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// GetModelInfo fetches detailed information about a specific model.
func (c *Client) GetModelInfo(ctx context.Context, modelID string) (*ModelInfo, error) {
	url := fmt.Sprintf("%s/models/%s", c.baseURL, modelID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var info ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}

// HealthCheck performs a simple GET request to the root path to verify connectivity.
func (c *Client) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health check returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
