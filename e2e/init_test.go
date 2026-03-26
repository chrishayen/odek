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

	out, code := run(t, dir, "init")
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

	// Verify toml has server config
	data, _ := os.ReadFile(filepath.Join(dir, "valkyrie.toml"))
	if !strings.Contains(string(data), "[server]") {
		t.Errorf("expected [server] in toml, got: %s", data)
	}
	if !strings.Contains(string(data), "http://localhost:7777") {
		t.Errorf("expected default server URL in toml, got: %s", data)
	}
}

func TestInitWithProjectName(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	out, code := run(t, dir, "init", "--project", "my-cool-app")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "my-cool-app") {
		t.Errorf("expected project name in output, got: %s", out)
	}
}

func TestInitWithCustomServer(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	out, code := run(t, dir, "init", "--server", "http://runes.example.com:9999")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "valkyrie.toml"))
	if !strings.Contains(string(data), "http://runes.example.com:9999") {
		t.Errorf("expected custom server URL in toml, got: %s", data)
	}
	_ = out
}

func TestInitFailsIfAlreadyInitialized(t *testing.T) {
	dir, err := os.MkdirTemp("", "valkyrie-init-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	_, code := run(t, dir, "init")
	if code != 0 {
		t.Fatal("first init should succeed")
	}

	out, code := run(t, dir, "init")
	if code == 0 {
		t.Error("expected non-zero exit on second init")
	}
	if !strings.Contains(out, "already") {
		t.Errorf("expected 'already' in error message, got: %s", out)
	}
}
