package e2e_test

import (
	"strings"
	"testing"
)

func TestCLIExportFeatureNotFound(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	_, code := run(t, dir, "features", "export", "nonexistent")
	if code == 0 {
		t.Error("expected non-zero exit for missing feature")
	}
}

func TestCLIExportFeatureNotHydrated(t *testing.T) {
	dir, cleanup := testEnv(t, "language = \"ts\"\n\n[agent]\nmock = true\n")
	defer cleanup()

	// Create a feature rune but don't hydrate it.
	run(t, dir, "runes", "create", "--name", "unhydrated", "--description", "Not hydrated", "--signature", "() -> void")

	out, code := run(t, dir, "features", "export", "unhydrated")
	if code == 0 {
		t.Errorf("expected non-zero exit for un-hydrated feature, got: %s", out)
	}
	if !strings.Contains(out, "unhydrated") {
		t.Errorf("expected 'unhydrated' in error, got: %s", out)
	}
}
