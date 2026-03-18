package e2e_test

import (
	"net/http"
	"strings"
	"testing"
)



func TestAuthRequired(t *testing.T) {
	base, cleanup := startServerWithToken(t, "test-secret")
	defer cleanup()

	// no token — should get 401
	resp, err := http.Get(base + "/runes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

func TestAuthWithValidToken(t *testing.T) {
	base, cleanup := startServerWithToken(t, "test-secret")
	defer cleanup()

	req, _ := http.NewRequest(http.MethodGet, base+"/runes", nil)
	req.Header.Set("Authorization", "Bearer test-secret")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200 with valid token, got %d", resp.StatusCode)
	}
}

func TestAuthWithInvalidToken(t *testing.T) {
	base, cleanup := startServerWithToken(t, "test-secret")
	defer cleanup()

	req, _ := http.NewRequest(http.MethodGet, base+"/runes", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 with wrong token, got %d", resp.StatusCode)
	}
}

func TestHealthPublicWithAuth(t *testing.T) {
	base, cleanup := startServerWithToken(t, "test-secret")
	defer cleanup()

	// /health is always public — no token needed
	resp, err := http.Get(base + "/health")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected /health to be public, got %d", resp.StatusCode)
	}
}

func TestAuthDisabled(t *testing.T) {
	base, cleanup := startServer(t, "") // helpers disable auth by default
	defer cleanup()

	// no token, should work fine
	resp, err := http.Get(base + "/runes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200 with auth disabled, got %d", resp.StatusCode)
	}
}

func TestServeFailsWithNoToken(t *testing.T) {
	// auth not disabled and no token set — should fail with helpful error
	out, code := runBinary(t, `[auth]
`, "serve")
	if code == 0 {
		t.Fatal("expected non-zero exit when auth not configured")
	}
	if !strings.Contains(out, "auth.token") {
		t.Errorf("expected helpful error about auth.token, got: %s", out)
	}
}
