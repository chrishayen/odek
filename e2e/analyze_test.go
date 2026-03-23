package e2e_test

import (
	"strings"
	"testing"
)

func TestAnalyzeDecomposes(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	out, code := run(t, dir, "runes", "analyze",
		"--requirements", "User login with email and password",
		"--yes")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "New runes") {
		t.Errorf("expected 'New runes' section in output, got: %s", out)
	}
	if !strings.Contains(out, "validate-email") {
		t.Errorf("expected 'validate-email' rune in output, got: %s", out)
	}
}

func TestAnalyzeCreatesRunes(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	out, code := run(t, dir, "runes", "analyze",
		"--requirements", "User login with email and password",
		"--yes")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, `created rune "auth/validate-email"`) {
		t.Errorf("expected auth/validate-email to be created, got: %s", out)
	}
	if !strings.Contains(out, `created rune "auth/hash-password"`) {
		t.Errorf("expected auth/hash-password to be created, got: %s", out)
	}
	if !strings.Contains(out, `created rune "auth/store-credentials"`) {
		t.Errorf("expected auth/store-credentials to be created, got: %s", out)
	}

	// Verify rune exists with rich fields
	listOut, _ := run(t, dir, "runes", "get", "auth/validate-email")
	if !strings.Contains(listOut, "validate-email") {
		t.Errorf("expected validate-email in registry, got: %s", listOut)
	}
	if !strings.Contains(listOut, "behavior") {
		t.Errorf("expected behavior field in rune, got: %s", listOut)
	}
}

func TestAnalyzeFindsExisting(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "auth/validate-email", "--description", "Validates email format")

	out, code := run(t, dir, "runes", "analyze",
		"--requirements", "User login with email and password",
		"--yes")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, `created rune "auth/hash-password"`) {
		t.Errorf("expected hash-password to be created, got: %s", out)
	}
}

func TestAnalyzeMissingRequirements(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	_, code := run(t, dir, "runes", "analyze")
	if code == 0 {
		t.Error("expected non-zero exit when --requirements is missing")
	}
}

func TestAnalyzeDefaultAgent(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	out, code := run(t, dir, "runes", "analyze",
		"--requirements", "Simple math operations",
		"--yes")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
}
