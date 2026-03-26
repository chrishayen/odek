package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/chrishayen/valkyrie/internal/server/store"
)

// Client talks to the valkyrie rune server API.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{},
	}
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("server request failed: %w", err)
	}
	return resp, nil
}

func decode[T any](resp *http.Response) (T, error) {
	var v T
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return v, fmt.Errorf("server error (%d): %s", resp.StatusCode, errResp.Error)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

// fqnToPath converts "net.http.parse_url" to "net/http/parse_url" for URL paths.
func fqnToPath(fqn string) string {
	return strings.ReplaceAll(fqn, ".", "/")
}

// Health checks the server is up.
func (c *Client) Health() error {
	resp, err := c.do("GET", "/api/health", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("server unhealthy: %d", resp.StatusCode)
	}
	return nil
}

// RunesList returns runes, optionally filtered.
func (c *Client) RunesList(project, namespace string) ([]store.Rune, error) {
	params := url.Values{}
	if project != "" {
		params.Set("project", project)
	}
	if namespace != "" {
		params.Set("namespace", namespace)
	}
	path := "/api/runes"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return decode[[]store.Rune](resp)
}

// RunesGet retrieves a single rune by FQN.
func (c *Client) RunesGet(fqn string) (*store.Rune, error) {
	resp, err := c.do("GET", "/api/runes/"+fqnToPath(fqn), nil)
	if err != nil {
		return nil, err
	}
	return decode[*store.Rune](resp)
}

// RunesCreate creates a new rune.
func (c *Client) RunesCreate(r store.Rune) (*store.Rune, error) {
	resp, err := c.do("POST", "/api/runes", r)
	if err != nil {
		return nil, err
	}
	return decode[*store.Rune](resp)
}

// RunesUpdate patches an existing rune.
func (c *Client) RunesUpdate(fqn string, patch map[string]any) (*store.Rune, error) {
	resp, err := c.do("PUT", "/api/runes/"+fqnToPath(fqn), patch)
	if err != nil {
		return nil, err
	}
	return decode[*store.Rune](resp)
}

// RunesDelete deletes a rune.
func (c *Client) RunesDelete(fqn string) error {
	resp, err := c.do("DELETE", "/api/runes/"+fqnToPath(fqn), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("delete failed: %d", resp.StatusCode)
	}
	return nil
}

// RunesSearch searches runes by query string.
func (c *Client) RunesSearch(query string) ([]store.Rune, error) {
	resp, err := c.do("POST", "/api/runes_search", map[string]string{"query": query})
	if err != nil {
		return nil, err
	}
	return decode[[]store.Rune](resp)
}

// RunesApprove marks a rune as approved.
func (c *Client) RunesApprove(fqn string) (*store.Rune, error) {
	resp, err := c.do("POST", "/api/runes_commit/"+fqnToPath(fqn), nil)
	if err != nil {
		return nil, err
	}
	return decode[*store.Rune](resp)
}

// ProjectsList returns all project names.
func (c *Client) ProjectsList() ([]string, error) {
	resp, err := c.do("GET", "/api/projects", nil)
	if err != nil {
		return nil, err
	}
	return decode[[]string](resp)
}

// ProjectRunes returns all runes for a project.
func (c *Client) ProjectRunes(name string) ([]store.Rune, error) {
	resp, err := c.do("GET", "/api/projects/"+name, nil)
	if err != nil {
		return nil, err
	}
	return decode[[]store.Rune](resp)
}

// RequirementsSubmit submits requirements for decomposition.
func (c *Client) RequirementsSubmit(project, requirements string) (map[string]string, error) {
	resp, err := c.do("POST", "/api/requirements", map[string]string{
		"project":      project,
		"requirements": requirements,
	})
	if err != nil {
		return nil, err
	}
	return decode[map[string]string](resp)
}

// RequirementsStatus checks the status of a requirements job.
func (c *Client) RequirementsStatus(id string) (map[string]string, error) {
	resp, err := c.do("GET", "/api/requirements/"+id, nil)
	if err != nil {
		return nil, err
	}
	return decode[map[string]string](resp)
}
