package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIExportFeature(t *testing.T) {
	dir, cleanup := testEnv(t, "language = \"ts\"\n\n[agent]\nmock = true\n")
	defer cleanup()

	// Create a top-level rune that qualifies as a feature.
	run(t, dir, "runes", "create", "--name", "greeter", "--description", "Greets a user by name", "--signature", "(name: string) -> string")

	// Manually mark it hydrated (mock hydrate generates Go, not TS).
	runeFile := filepath.Join(dir, "runes", "greeter", "1.0.0.md")
	data, err := os.ReadFile(runeFile)
	if err != nil {
		t.Fatalf("reading rune file: %v", err)
	}
	content := strings.Replace(string(data), "hydrated: false", "hydrated: true", 1)
	content = strings.Replace(content, "status: ", "status: stable", 1)
	if !strings.Contains(content, "status: stable") {
		content = strings.Replace(content, "hydrated: true", "status: stable\nhydrated: true", 1)
	}
	os.WriteFile(runeFile, []byte(content), 0644)

	// Create source files to mimic hydrated + composed output.
	srcDir := filepath.Join(dir, "src", "greeter")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "greeter.ts"), []byte(
		"export function greet(name: string): string { return `Hello, ${name}!`; }\n",
	), 0644)
	os.WriteFile(filepath.Join(srcDir, "greeter.test.ts"), []byte(
		"import { greet } from './greeter.ts';\nconsole.log(greet('world'));\n",
	), 0644)

	// Create dispatch framework file (normally done by compose).
	dispatchDir := filepath.Join(dir, "src", "dispatch")
	os.MkdirAll(dispatchDir, 0755)
	os.WriteFile(filepath.Join(dispatchDir, "dispatch.ts"), []byte(
		"export class Dispatcher {}\n",
	), 0644)

	// Run export.
	out, code := run(t, dir, "features", "export", "greeter")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}

	// Parse result.
	var result struct {
		FeatureName string   `json:"feature_name"`
		Version     string   `json:"version"`
		OutputDir   string   `json:"output_dir"`
		Files       []string `json:"files"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("parsing output: %v\nraw: %s", err, out)
	}

	if result.FeatureName != "greeter" {
		t.Errorf("feature_name = %q, want %q", result.FeatureName, "greeter")
	}
	if result.Version != "1.0.0" {
		t.Errorf("version = %q, want %q", result.Version, "1.0.0")
	}

	// Verify dist directory structure.
	distDir := filepath.Join(dir, "dist", "greeter")
	for _, f := range []string{"greeter.ts", "dispatch/dispatch.ts", "index.ts", "package.json"} {
		path := filepath.Join(distDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %q to exist in dist", f)
		}
	}

	// Verify test files excluded.
	testFile := filepath.Join(distDir, "greeter.test.ts")
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("test file should not exist in dist")
	}

	// Verify package.json.
	pkgData, err := os.ReadFile(filepath.Join(distDir, "package.json"))
	if err != nil {
		t.Fatalf("reading package.json: %v", err)
	}
	var pkg map[string]any
	json.Unmarshal(pkgData, &pkg)
	if pkg["name"] != "greeter" {
		t.Errorf("package name = %q, want %q", pkg["name"], "greeter")
	}
	if pkg["version"] != "1.0.0" {
		t.Errorf("package version = %q, want %q", pkg["version"], "1.0.0")
	}
	if pkg["type"] != "module" {
		t.Errorf("package type = %q, want %q", pkg["type"], "module")
	}
}

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
