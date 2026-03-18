package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func apiPost(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	return resp
}

func apiGet(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return resp
}

func decodeBody(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
}

// --- Tests ---

func TestAPIHealth(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	resp := apiGet(t, base+"/health")
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAPICreateRune(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	resp := apiPost(t, base+"/runes", map[string]any{
		"name":        "user-auth",
		"description": "Handles user authentication via JWT",
		"runtime":     "go@1.22",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var rune map[string]any
	decodeBody(t, resp, &rune)

	if rune["name"] != "user-auth" {
		t.Errorf("expected name=user-auth, got %v", rune["name"])
	}
	if rune["stage"] != "draft" {
		t.Errorf("expected stage=draft, got %v", rune["stage"])
	}
	if rune["version"] != "0.1.0" {
		t.Errorf("expected version=0.1.0, got %v", rune["version"])
	}
}

func TestAPICreateRuneValidation(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	// missing description
	resp := apiPost(t, base+"/runes", map[string]any{"name": "bad"})
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing description, got %d", resp.StatusCode)
	}

	// missing name
	resp = apiPost(t, base+"/runes", map[string]any{"description": "no name"})
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing name, got %d", resp.StatusCode)
	}

	// duplicate
	apiPost(t, base+"/runes", map[string]any{"name": "dup", "description": "first"})
	resp = apiPost(t, base+"/runes", map[string]any{"name": "dup", "description": "second"})
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for duplicate, got %d", resp.StatusCode)
	}
}

func TestAPIListRunes(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	apiPost(t, base+"/runes", map[string]any{"name": "rune-a", "description": "first"})
	apiPost(t, base+"/runes", map[string]any{"name": "rune-b", "description": "second"})

	resp := apiGet(t, base+"/runes")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var runes []map[string]any
	decodeBody(t, resp, &runes)
	if len(runes) != 2 {
		t.Errorf("expected 2 runes, got %d", len(runes))
	}
}

func TestAPIGetRune(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	apiPost(t, base+"/runes", map[string]any{"name": "my-rune", "description": "test rune"})

	resp := apiGet(t, base+"/runes/my-rune")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var rune map[string]any
	decodeBody(t, resp, &rune)
	if rune["name"] != "my-rune" {
		t.Errorf("expected my-rune, got %v", rune["name"])
	}
}

func TestAPIGetRuneNotFound(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	resp := apiGet(t, base+"/runes/does-not-exist")
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestAPIPromoteRune(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	apiPost(t, base+"/runes", map[string]any{"name": "promo-rune", "description": "to promote"})

	// draft → reviewed
	resp, _ := http.Post(base+"/runes/promo-rune/promote", "application/json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on first promote, got %d", resp.StatusCode)
	}
	var rune map[string]any
	decodeBody(t, resp, &rune)
	if rune["stage"] != "reviewed" {
		t.Errorf("expected reviewed, got %v", rune["stage"])
	}

	// reviewed → stable
	resp, _ = http.Post(base+"/runes/promo-rune/promote", "application/json", nil)
	decodeBody(t, resp, &rune)
	if rune["stage"] != "stable" {
		t.Errorf("expected stable, got %v", rune["stage"])
	}

	// stable → error
	resp, _ = http.Post(base+"/runes/promo-rune/promote", "application/json", nil)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 when already stable, got %d", resp.StatusCode)
	}
}

func TestAPIDeleteRune(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	apiPost(t, base+"/runes", map[string]any{"name": "delete-me", "description": "gone soon"})

	req, _ := http.NewRequest(http.MethodDelete, base+"/runes/delete-me", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 204 {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}

	resp = apiGet(t, base+"/runes/delete-me")
	if resp.StatusCode != 404 {
		t.Errorf("expected 404 after delete, got %d", resp.StatusCode)
	}
}

func TestAPIUpdateRune(t *testing.T) {
	base, cleanup := startServer(t, "")
	defer cleanup()

	apiPost(t, base+"/runes", map[string]any{"name": "update-me", "description": "original"})

	body, _ := json.Marshal(map[string]any{
		"name":        "update-me",
		"description": "updated description",
		"version":     "0.2.0",
		"stage":       "draft",
	})
	req, _ := http.NewRequest(http.MethodPut, base+"/runes/update-me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var rune map[string]any
	decodeBody(t, resp, &rune)
	if !strings.Contains(fmt.Sprint(rune["description"]), "updated") {
		t.Errorf("expected updated description, got %v", rune["description"])
	}
}
