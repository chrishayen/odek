package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chrishayen/valkyrie/internal/adaptor"
	"github.com/chrishayen/valkyrie/internal/server"
)

// mockAdaptor returns canned JSON responses based on prompt content.
type mockAdaptor struct {
	project string
}

func (m *mockAdaptor) Run(_ context.Context, prompt string) (string, error) {
	if strings.Contains(prompt, "decomposing software requirements") {
		return fmt.Sprintf(`[
			{"name": "validate_email", "description": "Validates email format per RFC 5322", "pure": true},
			{"name": "hash_password", "description": "Hashes a password with bcrypt", "pure": true},
			{"name": "create_session", "description": "Creates a session token for authenticated user", "pure": true}
		]`), nil
	}

	if strings.Contains(prompt, "classifying functions") {
		return fmt.Sprintf(`[
			{"name": "validate_email", "description": "Validates email format per RFC 5322", "pure": true, "fqn": "text.validate.email", "existing": false},
			{"name": "hash_password", "description": "Hashes a password with bcrypt", "pure": true, "fqn": "crypto.hash.bcrypt", "existing": false},
			{"name": "create_session", "description": "Creates a session token", "pure": true, "fqn": "%s.auth.create_session", "existing": false}
		]`, m.project), nil
	}

	if strings.Contains(prompt, "designing a rune") {
		// Extract the FQN from the prompt
		fqn := "unknown.rune"
		if idx := strings.Index(prompt, "FQN: "); idx != -1 {
			end := strings.Index(prompt[idx+5:], "\n")
			if end != -1 {
				fqn = prompt[idx+5 : idx+5+end]
			}
		}
		return fmt.Sprintf(`{
			"fqn": %q,
			"description": "Test rune for %s",
			"signature": "(input: string) -> result[string, error]",
			"behavior": "- Takes input string\n- Returns processed result\n- Errors on empty input",
			"positive_tests": ["valid input returns success", "unicode input handled correctly"],
			"negative_tests": ["empty string returns error", "null input returns error"]
		}`, fqn, fqn), nil
	}

	return "[]", nil
}

var _ adaptor.Adaptor = (*mockAdaptor)(nil)

// startMockServer starts a server in-process with the mock adaptor.
func startMockServer(t *testing.T, project string) (url string, cleanup func()) {
	t.Helper()

	dataDir, err := os.MkdirTemp("", "valkyrie-pipeline-*")
	if err != nil {
		t.Fatal(err)
	}

	mock := &mockAdaptor{project: project}
	s := server.New(dataDir, "test-token", mock)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go http.Serve(l, s)

	url = fmt.Sprintf("http://%s", l.Addr().String())
	cleanup = func() {
		l.Close()
		os.RemoveAll(dataDir)
	}
	return url, cleanup
}

func TestPipelineDecomposesRequirements(t *testing.T) {
	url, cleanup := startMockServer(t, "myapp")
	defer cleanup()

	// Submit requirements
	result, code := doJSON[map[string]string](t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "Build a user authentication system with email/password login and session management",
	}))
	if code != 202 {
		t.Fatalf("expected 202, got %d", code)
	}
	jobID := result["id"]
	if jobID == "" {
		t.Fatal("expected job ID")
	}

	// Poll until complete (mock is fast)
	var job map[string]any
	for i := 0; i < 50; i++ {
		job, code = doJSON[map[string]any](t, authReq(t, "GET", url+"/api/requirements/"+jobID, nil))
		if code != 200 {
			t.Fatalf("poll: expected 200, got %d", code)
		}
		status, _ := job["status"].(string)
		if status == "completed" || status == "failed" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	status, _ := job["status"].(string)
	if status != "completed" {
		errMsg, _ := job["error"].(string)
		t.Fatalf("expected completed, got %s (error: %s)", status, errMsg)
	}

	// Parse result
	resultRaw, _ := json.Marshal(job["result"])
	var pipelineResult struct {
		Proposals []struct {
			FQN      string `json:"fqn"`
			Existing bool   `json:"existing"`
			Spec     *struct {
				FQN         string   `json:"fqn"`
				Description string   `json:"description"`
				Signature   string   `json:"signature"`
				Behavior    string   `json:"behavior"`
				PosTests    []string `json:"positive_tests"`
				NegTests    []string `json:"negative_tests"`
				Status      string   `json:"status"`
			} `json:"spec"`
		} `json:"proposals"`
	}
	if err := json.Unmarshal(resultRaw, &pipelineResult); err != nil {
		t.Fatalf("parsing result: %v\nraw: %s", err, string(resultRaw))
	}

	if len(pipelineResult.Proposals) != 3 {
		t.Fatalf("expected 3 proposals, got %d", len(pipelineResult.Proposals))
	}

	// Check FQNs
	expectedFQNs := map[string]bool{
		"text.validate.email":     true,
		"crypto.hash.bcrypt":      true,
		"myapp.auth.create_session": true,
	}
	for _, p := range pipelineResult.Proposals {
		if !expectedFQNs[p.FQN] {
			t.Errorf("unexpected FQN: %s", p.FQN)
		}
		if p.Existing {
			t.Errorf("expected all new, got existing for %s", p.FQN)
		}
		if p.Spec == nil {
			t.Errorf("expected spec for %s", p.FQN)
			continue
		}
		if p.Spec.Status != "draft" {
			t.Errorf("expected draft status for %s, got %s", p.FQN, p.Spec.Status)
		}
		if len(p.Spec.PosTests) == 0 {
			t.Errorf("expected positive tests for %s", p.FQN)
		}
		if len(p.Spec.NegTests) == 0 {
			t.Errorf("expected negative tests for %s", p.FQN)
		}
	}
}

func TestPipelineCreatesRunesOnServer(t *testing.T) {
	url, cleanup := startMockServer(t, "myapp")
	defer cleanup()

	// Submit and wait
	result, _ := doJSON[map[string]string](t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "Build auth",
	}))
	jobID := result["id"]

	for i := 0; i < 50; i++ {
		job, _ := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/requirements/"+jobID, nil))
		status, _ := job["status"].(string)
		if status == "completed" || status == "failed" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Verify runes were persisted on the server
	runes, code := doJSON[[]map[string]any](t, authReq(t, "GET", url+"/api/runes", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(runes) != 3 {
		t.Fatalf("expected 3 runes persisted, got %d", len(runes))
	}

	// Check we can get each one individually
	for _, fqn := range []string{"text.validate.email", "crypto.hash.bcrypt", "myapp.auth.create_session"} {
		pathFQN := strings.ReplaceAll(fqn, ".", "/")
		rune, code := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/"+pathFQN, nil))
		if code != 200 {
			t.Errorf("expected 200 for %s, got %d", fqn, code)
			continue
		}
		if rune["fqn"] != fqn {
			t.Errorf("expected fqn %s, got %v", fqn, rune["fqn"])
		}
		if rune["status"] != "draft" {
			t.Errorf("expected draft for %s, got %v", fqn, rune["status"])
		}
	}

	// Project rune should have project field set
	rune, _ := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/myapp/auth/create_session", nil))
	if rune["project"] != "myapp" {
		t.Errorf("expected project myapp, got %v", rune["project"])
	}

	// Generic runes should not have project field (omitted or empty)
	rune, _ = doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/text/validate/email", nil))
	proj, _ := rune["project"].(string)
	if proj != "" {
		t.Errorf("expected empty project for generic rune, got %v", rune["project"])
	}
}

func TestPipelineApproveAfterDecompose(t *testing.T) {
	url, cleanup := startMockServer(t, "myapp")
	defer cleanup()

	// Submit, wait for completion
	result, _ := doJSON[map[string]string](t, authReq(t, "POST", url+"/api/requirements", map[string]string{
		"project":      "myapp",
		"requirements": "Build auth",
	}))
	jobID := result["id"]

	for i := 0; i < 50; i++ {
		job, _ := doJSON[map[string]any](t, authReq(t, "GET", url+"/api/requirements/"+jobID, nil))
		status, _ := job["status"].(string)
		if status == "completed" || status == "failed" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Approve one rune
	rune, code := doJSON[map[string]any](t, authReq(t, "POST", url+"/api/runes_commit/text/validate/email", nil))
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if rune["status"] != "approved" {
		t.Errorf("expected approved, got %v", rune["status"])
	}

	// Others should still be draft
	rune, _ = doJSON[map[string]any](t, authReq(t, "GET", url+"/api/runes/crypto/hash/bcrypt", nil))
	if rune["status"] != "draft" {
		t.Errorf("expected draft, got %v", rune["status"])
	}
}
