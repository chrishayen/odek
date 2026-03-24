package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHydrateHelloWorld(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "hello-world-single", "--description", "Returns the string Hello World when called", "--signature", "() -> string")

	out, code := run(t, dir, "runes", "hydrate", "hello-world-single")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "hello-world-single") {
		t.Errorf("expected rune name in hydrate output, got: %s", out)
	}
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage in hydrate output, got: %s", out)
	}

	// Verify rune is marked hydrated
	out, _ = run(t, dir, "runes", "get", "hello-world-single")
	if !strings.Contains(out, `"hydrated": true`) {
		t.Errorf("expected hydrated=true, got: %s", out)
	}
}

func TestHydrateGeneratesCodeFiles(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "hello-world-files", "--description", "Returns Hello World", "--signature", "() -> string")
	run(t, dir, "runes", "hydrate", "hello-world-files")

	codeDir := filepath.Join(dir, "runes", "hello-world-files")
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		t.Fatalf("code dir not created: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated files")
	}
}

func TestHydrateCoverageTracked(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "hello-world-coverage", "--description", "Returns Hello World", "--signature", "() -> string")
	run(t, dir, "runes", "hydrate", "hello-world-coverage")

	out, _ := run(t, dir, "runes", "get", "hello-world-coverage")
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage field on rune, got: %s", out)
	}
}

func TestHydrateDefaultAgent(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\ntype = \"mock\"\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "hello-default", "--description", "Returns Hello World", "--signature", "() -> string")

	out, code := run(t, dir, "runes", "hydrate", "hello-default")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "hello-default") {
		t.Errorf("expected rune name in output, got: %s", out)
	}
}
