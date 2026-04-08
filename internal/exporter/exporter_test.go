package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

// setupTestFeature creates a rune registry with a feature that has an external dependency.
//
//	greeter (structural, no signature)
//	  greeter.greet (leaf, depends on std.io.write_stdout)
//	std (structural)
//	  std.io.write_stdout (leaf, no deps)
func setupTestFeature(t *testing.T) (string, *runepkg.Store) {
	t.Helper()
	dir := t.TempDir()
	runesDir := filepath.Join(dir, "runes")
	srcDir := filepath.Join(dir, "src")

	rs := runepkg.NewStore(runesDir, srcDir)

	// Feature: greeter (structural top-level).
	rs.Create(runepkg.Rune{
		Name:        "greeter",
		Description: "Greets the user",
	})

	// greeter.greet — leaf rune with external dependency.
	rs.Create(runepkg.Rune{
		Name:         "greeter.greet",
		Description:  "Writes hello world to stdout",
		Signature:    "() -> result[void, string]",
		Status:       "stable",
		Dependencies: []string{"std.io.write_stdout"},
	})
	r, _ := rs.Get("greeter.greet")
	r.Hydrated = true
	r.Coverage = 100.0
	rs.Update(*r)

	// std (structural).
	rs.Create(runepkg.Rune{
		Name:        "std",
		Description: "Standard library",
	})

	// std.io.write_stdout — leaf rune, no deps.
	rs.Create(runepkg.Rune{
		Name:        "std.io.write_stdout",
		Description: "Writes a string to stdout",
		Signature:   "(message: string) -> result[void, string]",
		Status:      "stable",
	})
	ws, _ := rs.Get("std.io.write_stdout")
	ws.Hydrated = true
	ws.Coverage = 100.0
	rs.Update(*ws)

	// Source files for greeter.greet (DI style — accepts write_stdout as param).
	greetSrc := filepath.Join(srcDir, "greeter")
	os.MkdirAll(greetSrc, 0755)
	os.WriteFile(filepath.Join(greetSrc, "greet.ts"), []byte(
		"export function greet(write_stdout: (msg: string) => any) { return write_stdout('hello world'); }\n",
	), 0644)
	os.WriteFile(filepath.Join(greetSrc, "greet.test.ts"), []byte(
		"import { greet } from './greet.ts';\n",
	), 0644)

	// Source files for std.io.write_stdout.
	stdSrc := filepath.Join(srcDir, "std", "io")
	os.MkdirAll(stdSrc, 0755)
	os.WriteFile(filepath.Join(stdSrc, "write_stdout.ts"), []byte(
		"export function write_stdout(message: string) { process.stdout.write(message); return { ok: true, value: undefined }; }\n",
	), 0644)
	os.WriteFile(filepath.Join(stdSrc, "write_stdout.test.ts"), []byte(
		"import { write_stdout } from './write_stdout.ts';\n",
	), 0644)

	return dir, rs
}

func TestExportProducesCorrectStructure(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if result.FeatureName != "greeter" {
		t.Errorf("FeatureName = %q, want %q", result.FeatureName, "greeter")
	}
	if result.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", result.Version, "1.0.0")
	}

	for _, f := range []string{
		"greet.ts",
		"std/io/write_stdout.ts",
		"wiring.ts",
		"index.ts",
		"package.json",
	} {
		path := filepath.Join(result.OutputDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %q to exist", f)
		}
	}
}

func TestExportExcludesTestFiles(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	for _, f := range result.Files {
		if strings.HasSuffix(f, ".test.ts") {
			t.Errorf("test file %q should not be in export", f)
		}
	}
}

func TestExportIncludesTestFiles(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{IncludeTests: true})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	hasTest := false
	for _, f := range result.Files {
		if strings.HasSuffix(f, ".test.ts") {
			hasTest = true
			break
		}
	}
	if !hasTest {
		t.Error("expected test files in export with IncludeTests=true")
	}
}

func TestExportPackageJSON(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(result.OutputDir, "package.json"))
	if err != nil {
		t.Fatalf("reading package.json: %v", err)
	}

	var pkg map[string]any
	json.Unmarshal(data, &pkg)
	if pkg["name"] != "greeter" {
		t.Errorf("name = %q, want %q", pkg["name"], "greeter")
	}
	if pkg["version"] != "1.0.0" {
		t.Errorf("version = %q, want %q", pkg["version"], "1.0.0")
	}
}

func TestExportWiringImportsDep(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(result.OutputDir, "wiring.ts"))
	if err != nil {
		t.Fatalf("reading wiring.ts: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "write_stdout") {
		t.Errorf("wiring should reference write_stdout: %s", content)
	}
	if !strings.Contains(content, "import") {
		t.Errorf("wiring should have imports: %s", content)
	}
}

func TestExportIndexReexportsWiring(t *testing.T) {
	dir, rs := setupTestFeature(t)
	distDir := filepath.Join(dir, "dist")
	exp := New(rs, "ts")

	result, err := exp.Export("greeter", distDir, ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(result.OutputDir, "index.ts"))
	if err != nil {
		t.Fatalf("reading index.ts: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "export * from './wiring'") {
		t.Errorf("index.ts should re-export wiring: %s", content)
	}
}

func TestExportNoDepsNoWiring(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(filepath.Join(dir, "runes"), filepath.Join(dir, "src"))

	// Simple feature with no dependencies.
	rs.Create(runepkg.Rune{
		Name:        "simple",
		Description: "A simple function",
		Signature:   "() -> string",
		Status:      "stable",
	})
	r, _ := rs.Get("simple")
	r.Hydrated = true
	rs.Update(*r)

	// Source file.
	srcDir := filepath.Join(dir, "src", "simple")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "simple.ts"), []byte("export function simple() { return 'hi'; }\n"), 0644)

	exp := New(rs, "ts")
	result, err := exp.Export("simple", filepath.Join(dir, "dist"), ExportOptions{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Should have no wiring file.
	for _, f := range result.Files {
		if f == "wiring.ts" {
			t.Error("should not have wiring.ts when there are no deps")
		}
	}

	// Index should re-export source directly.
	data, _ := os.ReadFile(filepath.Join(result.OutputDir, "index.ts"))
	if !strings.Contains(string(data), "export * from './simple'") {
		t.Errorf("index should re-export simple: %s", string(data))
	}
}

func TestExportFeatureNotFound(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(filepath.Join(dir, "runes"), filepath.Join(dir, "src"))
	exp := New(rs, "ts")

	_, err := exp.Export("nonexistent", filepath.Join(dir, "dist"), ExportOptions{})
	if err == nil {
		t.Fatal("expected error for nonexistent feature")
	}
}

func TestExportFeatureUnhydratedLeaf(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(filepath.Join(dir, "runes"), filepath.Join(dir, "src"))

	rs.Create(runepkg.Rune{
		Name:        "unhydrated",
		Description: "Has unhydrated children",
	})
	rs.Create(runepkg.Rune{
		Name:        "unhydrated.child",
		Description: "Not yet hydrated",
		Signature:   "() -> void",
		Status:      "stable",
	})

	exp := New(rs, "ts")
	_, err := exp.Export("unhydrated", filepath.Join(dir, "dist"), ExportOptions{})
	if err == nil {
		t.Fatal("expected error for un-hydrated leaf rune")
	}
	if !strings.Contains(err.Error(), "unhydrated.child") {
		t.Errorf("error should name the unhydrated rune: %v", err)
	}
}
