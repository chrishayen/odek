package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHydrateHelloWorld(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test.hello_world_single", "--description", "Returns the string Hello World when called", "--signature", "() -> string")

	out, code := run(t, dir, "runes", "hydrate", "test.hello_world_single")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "test.hello_world_single") {
		t.Errorf("expected rune name in hydrate output, got: %s", out)
	}
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage in hydrate output, got: %s", out)
	}

	// Verify rune is marked hydrated
	out, _ = run(t, dir, "runes", "get", "test.hello_world_single")
	if !strings.Contains(out, `"hydrated": true`) {
		t.Errorf("expected hydrated=true, got: %s", out)
	}
}

func TestHydrateGeneratesCodeFiles(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test.hello_world_files", "--description", "Returns Hello World", "--signature", "() -> string")
	run(t, dir, "runes", "hydrate", "test.hello_world_files")

	codeDir := filepath.Join(dir, "src", "test", "hello_world_files")
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		t.Fatalf("code dir not created: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated files")
	}
}

func TestHydrateCoverageTracked(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test.hello_world_coverage", "--description", "Returns Hello World", "--signature", "() -> string")
	run(t, dir, "runes", "hydrate", "test.hello_world_coverage")

	out, _ := run(t, dir, "runes", "get", "test.hello_world_coverage")
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage field on rune, got: %s", out)
	}
}

func TestHydrateDefaultAgent(t *testing.T) {
	dir, cleanup := testEnv(t, "[agent]\nmock = true\n")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test.hello_default", "--description", "Returns Hello World", "--signature", "() -> string")

	out, code := run(t, dir, "runes", "hydrate", "test.hello_default")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "test.hello_default") {
		t.Errorf("expected rune name in output, got: %s", out)
	}
}
