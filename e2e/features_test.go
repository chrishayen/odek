package e2e_test

import (
	"strings"
	"testing"
)

func TestCLIFeaturesList(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	// Create a top-level rune with a signature (becomes a feature when not draft).
	run(t, dir, "runes", "create", "--name", "auth", "--description", "Authentication feature", "--signature", "(email: string, password: string) -> bool")
	// Update status to stable so it appears as a feature.
	run(t, dir, "runes", "update", "auth", "--version", "1.0.0")

	out, code := run(t, dir, "features", "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected auth in output, got: %s", out)
	}
}

func TestCLIFeaturesGet(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "auth", "--description", "Authentication feature", "--signature", "(email: string, password: string) -> bool")

	out, code := run(t, dir, "features", "get", "auth")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected auth in output, got: %s", out)
	}
}

func TestCLIFeaturesGetNotFound(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	_, code := run(t, dir, "features", "get", "nosuchfeature")
	if code == 0 {
		t.Error("expected non-zero exit for missing feature")
	}
}
