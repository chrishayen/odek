package codegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractFiles(t *testing.T) {
	dir := t.TempDir()
	output := `=== FILE: main.go ===
package main

func Run() string { return "hello" }
=== END FILE ===

=== FILE: main_test.go ===
package main

import "testing"

func TestRun(t *testing.T) {}
=== END FILE ===
`
	if err := ExtractFiles(dir, output); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != "package main\n\nfunc Run() string { return \"hello\" }\n" {
		t.Errorf("main.go content = %q", got)
	}

	if _, err := os.Stat(filepath.Join(dir, "main_test.go")); err != nil {
		t.Errorf("main_test.go not created: %v", err)
	}
}

func TestExtractFilesNoBlocks(t *testing.T) {
	dir := t.TempDir()
	output := "package main\n\nfunc main() {}\n"

	if err := ExtractFiles(dir, output); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != output {
		t.Errorf("fallback content mismatch")
	}
}

func TestExtractFilesSubdir(t *testing.T) {
	dir := t.TempDir()
	output := `=== FILE: pkg/util.go ===
package pkg
=== END FILE ===
`
	if err := ExtractFiles(dir, output); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "pkg", "util.go")); err != nil {
		t.Errorf("subdirectory file not created: %v", err)
	}
}

func TestExtractFilesFlat(t *testing.T) {
	dir := t.TempDir()
	output := `=== FILE: src/auth/login.ts ===
export function login() {}
=== END FILE ===

=== FILE: src/auth/login.test.ts ===
import { login } from './login';
=== END FILE ===

=== FILE: package.json ===
{"name": "junk"}
=== END FILE ===

=== FILE: tsconfig.json ===
{}
=== END FILE ===
`
	if err := ExtractFilesFlat(dir, output); err != nil {
		t.Fatal(err)
	}

	// Source files should be written flat (no src/auth/ prefix)
	data, err := os.ReadFile(filepath.Join(dir, "login.ts"))
	if err != nil {
		t.Fatal("login.ts not created:", err)
	}
	if got := string(data); got != "export function login() {}\n" {
		t.Errorf("login.ts content = %q", got)
	}

	if _, err := os.Stat(filepath.Join(dir, "login.test.ts")); err != nil {
		t.Error("login.test.ts not created:", err)
	}

	// Config files should be skipped
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		t.Error("package.json should have been skipped")
	}
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); err == nil {
		t.Error("tsconfig.json should have been skipped")
	}

	// No nested directories should be created
	if _, err := os.Stat(filepath.Join(dir, "src")); err == nil {
		t.Error("nested src/ directory should not exist")
	}
}

func TestParseCoverage(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"ok  \tpkg\t0.001s\tcoverage: 85.7% of statements", 85.7},
		{"coverage: 100.0% of statements", 100.0},
		{"no coverage info here", -1},
		{"", -1},
	}
	for _, tt := range tests {
		got := ParseCoverage(tt.input)
		if got != tt.want {
			t.Errorf("ParseCoverage(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
