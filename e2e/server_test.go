package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func authReq(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()
	var req *http.Request
	if body != nil {
		data, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", "Bearer test-token")
	return req
}

func doJSON[T any](t *testing.T, req *http.Request) (T, int) {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var v T
	json.NewDecoder(resp.Body).Decode(&v)
	return v, resp.StatusCode
}

func doStatus(t *testing.T, req *http.Request) int {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

func createRune(t *testing.T, url string, r map[string]any) int {
	t.Helper()
	return doStatus(t, authReq(t, "POST", url+"/api/runes", r))
}

// --- Health ---

func TestServerHealthCheck(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	resp, err := http.Get(url + "/api/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

// --- Auth ---

func TestServerAuthRequired(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	// No token
	resp, _ := http.Get(url + "/api/runes")
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 without token, got %d", resp.StatusCode)
	}

	// Wrong token
	req, _ := http.NewRequest("GET", url+"/api/runes", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 with wrong token, got %d", resp.StatusCode)
	}

	// Correct token
	req, _ = http.NewRequest("GET", url+"/api/runes", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 with correct token, got %d", resp.StatusCode)
	}
}

func TestServerAuthOnAllEndpoints(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/runes"},
		{"POST", "/api/runes"},
		{"GET", "/api/runes/text/validate/email"},
		{"PUT", "/api/runes/text/validate/email"},
		{"DELETE", "/api/runes/text/validate/email"},
		{"POST", "/api/runes_search"},
		{"POST", "/api/runes_commit/text/validate/email"},
		{"GET", "/api/projects"},
		{"GET", "/api/projects/myapp"},
		{"POST", "/api/requirements"},
		{"GET", "/api/requirements/123"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, url+ep.path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", ep.method, ep.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != 401 {
			t.Errorf("%s %s: expected 401, got %d", ep.method, ep.path, resp.StatusCode)
		}
	}
}

// --- Rune CRUD ---

func TestServerRuneCreate(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	rune, code := doJSON[map[string]any](t, authReq(t, "POST", url+"/api/runes", map[string]string{
		"fqn":         "text.validate.email",
		"description": "Validates email format",
		"signature":   "(email: string) -> bool",
	}))
	if code != 201 {
		t.Fatalf("expected 201, got %d", code)
	}
	if rune["fqn"] != "text.validate.email" {
		t.Errorf("fqn: got %v", rune["fqn"])
	}
	if rune["status"] != "draft" {
		t.Errorf("status: expected draft, got %v", rune["status"])
	}
	if rune["version"] != "0.1.0" {
		t.Errorf("version: expected 0.1.0, got %v", rune["version"])
	}
}

func TestServerRuneCreateDuplicate(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	rune := map[string]any{
		"fqn":         "text.validate.email",
		"description": "Validates email format",
		"signature":   "(email: string) -> bool",
	}
	code := createRune(t, url, rune)
	if code != 201 {
		t.Fatalf("first create: expected 201, got %d", code)
	}

	code = createRune(t, url, rune)
	if code != 400 {
		t.Fatalf("duplicate create: expected 400, got %d", code)
	}
}

func TestServerRuneCreateValidation(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	cases := []struct {
		name string
		body map[string]any
	}{
		{"missing fqn", map[string]any{"description": "d", "signature": "s"}},
		{"missing description", map[string]any{"fqn": "a.b", "signature": "s"}},
		{"missing signature", map[string]any{"fqn": "a.b", "description": "d"}},
		{"single segment fqn", map[string]any{"fqn": "email", "description": "d", "signature": "s"}},
		{"empty segment in fqn", map[string]any{"fqn": "a..b", "description": "d", "signature": "s"}},
		{"slash in fqn", map[string]any{"fqn": "a/b.c", "description": "d", "signature": "s"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code := createRune(t, url, tc.body)
			if code != 400 {
				t.Errorf("expected 400, got %d", code)
			}
		})
	}
}

func TestServerRuneGet(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn":         "crypto.hash.sha256",
		"description": "Computes SHA-256 hash",
		"signature":   "(data: bytes) -> string",
		"behavior":    "Returns hex-encoded SHA-256 digest",
	})

	rune, code := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/crypto/hash/sha256", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if rune["fqn"] != "crypto.hash.sha256" {
		t.Errorf("fqn: got %v", rune["fqn"])
	}
	if rune["behavior"] != "Returns hex-encoded SHA-256 digest" {
		t.Errorf("behavior: got %v", rune["behavior"])
	}
}

func TestServerRuneGetNotFound(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "GET", url+"/api/runes/does/not/exist", nil))
	if code != 404 {
		t.Fatalf("expected 404, got %d", code)
	}
}

func TestServerRuneUpdate(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn":         "text.validate.email",
		"description": "Validates email format",
		"signature":   "(email: string) -> bool",
	})

	updated, code := doJSON[map[string]any](t, authReq(t, "PUT", url+"/api/runes/text/validate/email", map[string]any{
		"description": "Validates email address per RFC 5322",
		"version":     "0.2.0",
		"behavior":    "Checks local-part and domain",
	}))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if updated["description"] != "Validates email address per RFC 5322" {
		t.Errorf("description not updated: %v", updated["description"])
	}
	if updated["version"] != "0.2.0" {
		t.Errorf("version not updated: %v", updated["version"])
	}
	if updated["behavior"] != "Checks local-part and domain" {
		t.Errorf("behavior not updated: %v", updated["behavior"])
	}
	// Signature should be unchanged
	if updated["signature"] != "(email: string) -> bool" {
		t.Errorf("signature should be unchanged: %v", updated["signature"])
	}
}

func TestServerRuneUpdateNotFound(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "PUT", url+"/api/runes/does/not/exist", map[string]any{"description": "x"}))
	if code != 404 {
		t.Fatalf("expected 404, got %d", code)
	}
}

func TestServerRuneDelete(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn":         "text.validate.email",
		"description": "Validates email format",
		"signature":   "(email: string) -> bool",
	})

	code := doStatus(t, authReq(t, "DELETE", url+"/api/runes/text/validate/email", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}

	// Verify gone
	code = doStatus(t, authReq(t, "GET", url+"/api/runes/text/validate/email", nil))
	if code != 404 {
		t.Fatalf("expected 404 after delete, got %d", code)
	}
}

func TestServerRuneDeleteNotFound(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "DELETE", url+"/api/runes/does/not/exist", nil))
	if code != 404 {
		t.Fatalf("expected 404, got %d", code)
	}
}

// --- List & Filter ---

func TestServerRuneListEmpty(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 0 {
		t.Errorf("expected empty list, got %d", len(runes))
	}
}

func TestServerRuneListAll(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	for _, fqn := range []string{
		"text.validate.email",
		"crypto.hash.sha256",
		"net.http.parse_url",
	} {
		createRune(t, url, map[string]any{
			"fqn": fqn, "description": fqn, "signature": "() -> string",
		})
	}

	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 3 {
		t.Errorf("expected 3 runes, got %d", len(runes))
	}
}

func TestServerRuneListFilterByNamespace(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	for _, fqn := range []string{
		"text.validate.email",
		"text.validate.phone",
		"text.format.title_case",
		"crypto.hash.sha256",
	} {
		createRune(t, url, map[string]any{
			"fqn": fqn, "description": fqn, "signature": "() -> string",
		})
	}

	// Filter by text.validate
	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?namespace=text.validate", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 2 {
		t.Errorf("expected 2 runes in text.validate, got %d", len(runes))
	}

	// Filter by text (broader)
	runes, code = doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?namespace=text", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 3 {
		t.Errorf("expected 3 runes in text.*, got %d", len(runes))
	}

	// Filter by crypto
	runes, code = doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?namespace=crypto", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 1 {
		t.Errorf("expected 1 rune in crypto.*, got %d", len(runes))
	}
}

func TestServerRuneListFilterByProject(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	// Generic rune
	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "generic", "signature": "() -> bool",
	})

	// Project-specific rune
	data, _ := json.Marshal(map[string]any{
		"fqn": "myapp.auth.validate_token", "description": "project rune", "signature": "() -> bool", "project": "myapp",
	})
	req, _ := http.NewRequest("POST", url+"/api/runes", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Filter by project
	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?project=myapp", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 1 {
		t.Errorf("expected 1 project rune, got %d", len(runes))
	}
	if len(runes) > 0 && runes[0]["fqn"] != "myapp.auth.validate_token" {
		t.Errorf("expected myapp rune, got %v", runes[0]["fqn"])
	}
}

// --- Commit (approve) ---

func TestServerRuneCommit(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "validates email", "signature": "() -> bool",
	})

	// Verify starts as draft
	rune, _ := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/text/validate/email", nil))
	if rune["status"] != "draft" {
		t.Fatalf("expected draft, got %v", rune["status"])
	}

	// Commit
	rune, code := doJSON[map[string]any](t, authReq(t, "POST", url+"/api/runes_commit/text/validate/email", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if rune["status"] != "approved" {
		t.Errorf("expected approved, got %v", rune["status"])
	}

	// Verify persisted
	rune, _ = doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/text/validate/email", nil))
	if rune["status"] != "approved" {
		t.Errorf("expected approved after re-read, got %v", rune["status"])
	}
}

func TestServerRuneCommitNotFound(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "POST", url+"/api/runes_commit/does/not/exist", nil))
	if code != 404 {
		t.Fatalf("expected 404, got %d", code)
	}
}

// --- Search ---

func TestServerSearchByFQN(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	for _, fqn := range []string{
		"text.validate.email",
		"text.validate.phone",
		"crypto.hash.sha256",
	} {
		createRune(t, url, map[string]any{
			"fqn": fqn, "description": fmt.Sprintf("Does %s things", fqn), "signature": "() -> string",
		})
	}

	results, code := doJSON[[]map[string]any](t, authReq(t, "POST", url+"/api/runes_search", map[string]string{"query": "sha256"}))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0]["fqn"] != "crypto.hash.sha256" {
		t.Errorf("wrong result: %v", results[0]["fqn"])
	}
}

func TestServerSearchByDescription(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "Validates email addresses per RFC 5322", "signature": "() -> bool",
	})
	createRune(t, url, map[string]any{
		"fqn": "crypto.hash.sha256", "description": "Computes SHA-256 digest", "signature": "() -> string",
	})

	results, code := doJSON[[]map[string]any](t, authReq(t, "POST", url+"/api/runes_search", map[string]string{"query": "RFC"}))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0]["fqn"] != "text.validate.email" {
		t.Errorf("wrong result: %v", results[0]["fqn"])
	}
}

func TestServerSearchCaseInsensitive(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "Validates EMAIL format", "signature": "() -> bool",
	})

	results, code := doJSON[[]map[string]any](t, authReq(t, "POST", url+"/api/runes_search", map[string]string{"query": "email"}))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}
}

func TestServerSearchNoResults(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "validates email", "signature": "() -> bool",
	})

	results, code := doJSON[[]map[string]any](t, authReq(t, "POST", url+"/api/runes_search", map[string]string{"query": "kubernetes"}))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// --- Deep nesting ---

func TestServerDeepNamespaceNesting(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	deep := "net.http.middleware.cors.parse_headers"
	createRune(t, url, map[string]any{
		"fqn": deep, "description": "Parses CORS headers", "signature": "(headers: map[string, string]) -> CorsConfig",
	})

	rune, code := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/net/http/middleware/cors/parse_headers", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if rune["fqn"] != deep {
		t.Errorf("expected %s, got %v", deep, rune["fqn"])
	}

	// Namespace filter at various levels
	for _, ns := range []string{"net", "net.http", "net.http.middleware", "net.http.middleware.cors"} {
		runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?namespace="+ns, nil))
		if code != 200 {
			t.Fatalf("expected 200 for namespace %s, got %d", ns, code)
		}
		if len(runes) != 1 {
			t.Errorf("expected 1 rune for namespace %s, got %d", ns, len(runes))
		}
	}
}

// --- Projects ---

func TestServerProjectsList(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	// Empty initially
	projects, code := doJSON[[]string](t, authReq(t, "GET", url+"/api/projects", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(projects))
	}
}

func TestServerProjectRunes(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	// Create runes under a project namespace
	for _, fqn := range []string{
		"myapp.auth.validate_token",
		"myapp.auth.hash_password",
		"myapp.payment.calculate_total",
	} {
		createRune(t, url, map[string]any{
			"fqn": fqn, "description": fqn, "signature": "() -> string",
		})
	}
	// And a generic rune
	createRune(t, url, map[string]any{
		"fqn": "text.validate.email", "description": "generic", "signature": "() -> bool",
	})

	// Get project runes via namespace
	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/projects/myapp", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 3 {
		t.Errorf("expected 3 myapp runes, got %d", len(runes))
	}

	// Subnamespace
	runes, code = doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes?namespace=myapp.auth", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 2 {
		t.Errorf("expected 2 myapp.auth runes, got %d", len(runes))
	}
}

// --- Requirements (stub) ---

func TestServerRequirementsSubmit(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	result, code := doJSON[map[string]string](t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "Build a user authentication system with email/password login",
	}))
	if code != 202 {
		t.Fatalf("expected 202, got %d", code)
	}
	if result["id"] == "" {
		t.Error("expected job id")
	}
	if result["status"] != "pending" {
		t.Errorf("expected pending, got %v", result["status"])
	}
}

func TestServerRequirementsSubmitEmpty(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "",
	}))
	if code != 400 {
		t.Fatalf("expected 400 for empty requirements, got %d", code)
	}
}

func TestServerRequirementsStatusNotFound(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	code := doStatus(t, authReq(t, "GET", url+"/api/requirements/nonexistent-id", nil))
	if code != 404 {
		t.Fatalf("expected 404 for unknown job, got %d", code)
	}
}

func TestServerRequirementsStatusAfterSubmit(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	// Submit requirements — pipeline will fail (fake API key) but job should be created
	result, code := doJSON[map[string]string](t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "Build a user auth system",
	}))
	if code != 202 {
		t.Fatalf("expected 202, got %d", code)
	}
	jobID := result["id"]

	// Should be able to get the job status
	var job map[string]any
	job, code = doJSON[map[string]any](t, authReq(t, "GET", url+"/api/requirements/"+jobID, nil))
	if code != 200 {
		t.Fatalf("expected 200 for known job, got %d", code)
	}
	if job["id"] != jobID {
		t.Errorf("expected id %s, got %v", jobID, job["id"])
	}
	// Status should be pending, running, completed, or failed
	status, _ := job["status"].(string)
	validStatuses := map[string]bool{"pending": true, "running": true, "completed": true, "failed": true}
	if !validStatuses[status] {
		t.Errorf("unexpected status: %v", status)
	}
}

// --- Rune spec round-trip (tests + behavior preserved) ---

func TestServerRuneFullSpecRoundTrip(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	original := map[string]any{
		"fqn":         "text.validate.email",
		"description": "Validates an email address format per RFC 5322",
		"signature":   "(email: string) -> result[bool, error]",
		"behavior":    "- Input: email string\n- Output: true if valid\n- Empty string returns error",
		"positive_tests": []string{
			"user@example.com returns true",
			"name+tag@domain.co.uk returns true",
		},
		"negative_tests": []string{
			"empty string returns error",
			"@nodomain returns false",
			"spaces in@email returns false",
		},
		"version": "1.0.0",
		"project": "myapp",
	}

	code := createRune(t, url, original)
	if code != 201 {
		t.Fatalf("create: expected 201, got %d", code)
	}

	rune, code := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/text/validate/email", nil))
	if code != 200 {
		t.Fatalf("get: expected 200, got %d", code)
	}

	if rune["fqn"] != original["fqn"] {
		t.Errorf("fqn: %v != %v", rune["fqn"], original["fqn"])
	}
	if rune["description"] != original["description"] {
		t.Errorf("description mismatch: got %v", rune["description"])
	}
	if rune["signature"] != original["signature"] {
		t.Errorf("signature mismatch: got %v", rune["signature"])
	}
	if rune["behavior"] != original["behavior"] {
		t.Errorf("behavior mismatch: got %v", rune["behavior"])
	}
	if rune["version"] != original["version"] {
		t.Errorf("version: %v != %v", rune["version"], original["version"])
	}
	if rune["project"] != original["project"] {
		t.Errorf("project: %v != %v", rune["project"], original["project"])
	}

	posTests, ok := rune["positive_tests"].([]any)
	if !ok {
		t.Fatalf("positive_tests: expected array, got %T", rune["positive_tests"])
	}
	if len(posTests) != 2 {
		t.Errorf("positive_tests: expected 2, got %d", len(posTests))
	}

	negTests, ok := rune["negative_tests"].([]any)
	if !ok {
		t.Fatalf("negative_tests: expected array, got %T", rune["negative_tests"])
	}
	if len(negTests) != 3 {
		t.Errorf("negative_tests: expected 3, got %d", len(negTests))
	}
}

// --- Invalid JSON ---

func TestServerInvalidJSON(t *testing.T) {
	url, cleanup := startServer(t)
	defer cleanup()

	req, _ := http.NewRequest("POST", url+"/api/runes", bytes.NewReader([]byte("not json")))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for bad JSON, got %d", resp.StatusCode)
	}
}
