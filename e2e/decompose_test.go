package e2e_test

import (
	"strings"
	"testing"
)

func TestDecomposeProducesTree(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	out, code := run(t, dir, "runes", "decompose",
		"--yes", "User login with email and password")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "New runes") {
		t.Errorf("expected 'New runes' section in output, got: %s", out)
	}
	if !strings.Contains(out, "validate_email") {
		t.Errorf("expected 'validate_email' rune in output, got: %s", out)
	}
}

func TestDecomposeCreatesRunes(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	out, code := run(t, dir, "runes", "decompose",
		"--yes", "User login with email and password")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, `created rune "std.auth.validate_email"`) {
		t.Errorf("expected std.auth.validate_email to be created, got: %s", out)
	}
	if !strings.Contains(out, `created rune "std.auth.hash_password"`) {
		t.Errorf("expected std.auth.hash_password to be created, got: %s", out)
	}

	// Verify rune exists
	listOut, _ := run(t, dir, "runes", "get", "std.auth.validate_email")
	if !strings.Contains(listOut, "validate_email") {
		t.Errorf("expected validate_email in registry, got: %s", listOut)
	}
}

func TestDecomposeMissingRequirements(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	_, code := run(t, dir, "runes", "decompose")
	if code == 0 {
		t.Error("expected non-zero exit when no args provided")
	}
}
