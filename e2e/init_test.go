package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCreatesProject(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	out, code := run(t, dir, "init", "go")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}

	// Check valkyrie.toml was created
	if _, err := os.Stat(filepath.Join(dir, "valkyrie.toml")); err != nil {
		t.Errorf("expected valkyrie.toml to exist: %v", err)
	}

	// Check .mcp.json was created
	if _, err := os.Stat(filepath.Join(dir, ".mcp.json")); err != nil {
		t.Errorf("expected .mcp.json to exist: %v", err)
	}

	// Check CLAUDE.md was created
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Errorf("expected CLAUDE.md to exist: %v", err)
	}

	// Check .claude/settings.json was created
	if _, err := os.Stat(filepath.Join(dir, ".claude", "settings.json")); err != nil {
		t.Errorf("expected .claude/settings.json to exist: %v", err)
	}

	// Project name should default to directory name
	dirName := filepath.Base(dir)
	if !strings.Contains(out, dirName) {
		t.Errorf("expected project name %q in output, got: %s", dirName, out)
	}
}

func TestInitWithProjectName(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	out, code := run(t, dir, "init", "go", "--project", "my-cool-app")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "my-cool-app") {
		t.Errorf("expected project name in output, got: %s", out)
	}
}

func TestInitFailsIfAlreadyInitialized(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// First init succeeds
	_, code := run(t, dir, "init", "go")
	if code != 0 {
		t.Fatal("first init should succeed")
	}

	// Second init fails
	out, code := run(t, dir, "init", "go")
	if code == 0 {
		t.Error("expected non-zero exit on second init")
	}
	if !strings.Contains(out, "already") {
		t.Errorf("expected 'already' in error message, got: %s", out)
	}
}
