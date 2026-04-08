package hydrator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/chrishayen/odek/internal/llm"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

func TestBuildPromptWithDeps(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(filepath.Join(dir, "runes"), filepath.Join(dir, "src"))

	// Create the dependency rune.
	rs.Create(runepkg.Rune{
		Name:        "std.io.write_stdout",
		Description: "Writes a string to stdout",
		Signature:   "(message: string) -> result[void, string]",
	})

	// Create the rune that depends on it.
	rs.Create(runepkg.Rune{
		Name:         "greeter.greet",
		Description:  "Writes hello world to stdout",
		Signature:    "() -> result[void, string]",
		Dependencies: []string{"std.io.write_stdout"},
	})

	client := llm.New("", "", true, "", "", 0)
	h := New(rs, client, "ts", nil)

	spec, err := h.GetHydrationSpec("greeter.greet")
	if err != nil {
		t.Fatalf("GetHydrationSpec failed: %v", err)
	}

	if !strings.Contains(spec.Prompt, "write_stdout") {
		t.Errorf("prompt should mention dependency write_stdout:\n%s", spec.Prompt)
	}
	if !strings.Contains(spec.Prompt, "dependencies as parameters") {
		t.Errorf("prompt should mention DI:\n%s", spec.Prompt)
	}
	if strings.Contains(spec.Prompt, "isolated from other runes") {
		t.Errorf("prompt should NOT say 'isolated from other runes' when deps exist:\n%s", spec.Prompt)
	}
}

func TestBuildPromptWithoutDeps(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(filepath.Join(dir, "runes"), filepath.Join(dir, "src"))

	rs.Create(runepkg.Rune{
		Name:        "std.io.write_stdout",
		Description: "Writes a string to stdout",
		Signature:   "(message: string) -> result[void, string]",
	})

	client := llm.New("", "", true, "", "", 0)
	h := New(rs, client, "ts", nil)

	spec, err := h.GetHydrationSpec("std.io.write_stdout")
	if err != nil {
		t.Fatalf("GetHydrationSpec failed: %v", err)
	}

	if strings.Contains(spec.Prompt, "dependencies as parameters") {
		t.Errorf("prompt should NOT mention DI for a leaf rune:\n%s", spec.Prompt)
	}
	if !strings.Contains(spec.Prompt, "no dependencies") {
		t.Errorf("prompt should say 'no dependencies':\n%s", spec.Prompt)
	}
}
