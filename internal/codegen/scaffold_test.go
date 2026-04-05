package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldFiles_TS(t *testing.T) {
	dir := t.TempDir()
	if err := ScaffoldFiles(dir, "login", ".ts"); err != nil {
		t.Fatal(err)
	}
	assertFileExists(t, filepath.Join(dir, "login.ts"))
	assertFileExists(t, filepath.Join(dir, "login.test.ts"))
}

func TestScaffoldFiles_Go(t *testing.T) {
	dir := t.TempDir()
	if err := ScaffoldFiles(dir, "login", ".go"); err != nil {
		t.Fatal(err)
	}
	assertFileExists(t, filepath.Join(dir, "login.go"))
	assertFileExists(t, filepath.Join(dir, "login_test.go"))
}

func TestScaffoldFiles_Py(t *testing.T) {
	dir := t.TempDir()
	if err := ScaffoldFiles(dir, "login", ".py"); err != nil {
		t.Fatal(err)
	}
	assertFileExists(t, filepath.Join(dir, "login.py"))
	assertFileExists(t, filepath.Join(dir, "test_login.py"))
}

func TestScaffoldFiles_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b")
	if err := ScaffoldFiles(dir, "foo", ".ts"); err != nil {
		t.Fatal(err)
	}
	assertFileExists(t, filepath.Join(dir, "foo.ts"))
}

func TestScaffoldFiles_Idempotent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "login.ts")
	if err := os.WriteFile(src, []byte("// existing"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ScaffoldFiles(dir, "login", ".ts"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(src)
	if string(data) != "// existing" {
		t.Errorf("ScaffoldFiles overwrote existing file, got %q", string(data))
	}
}

func TestScaffoldOnlyLeaves_Integration(t *testing.T) {
	// Reproduces the full decompose→scaffold flow.
	// Runes: hello_world (package), hello_world.greet (leaf),
	//        std (package), std.io (package), std.io.write_stdout (leaf)
	// Only leaves should get files. Packages should NOT.

	srcDir := t.TempDir()

	type testRune struct {
		name string
		leaf bool
	}
	runes := []testRune{
		{"hello_world", false},
		{"hello_world.greet", true},
		{"std", false},
		{"std.io", false},
		{"std.io.write_stdout", true},
	}

	allNames := make([]string, len(runes))
	for i, r := range runes {
		allNames[i] = r.name
	}

	ext := ".ts"
	for _, r := range runes {
		if !isLeaf(r.name, allNames) {
			continue
		}
		parts := strings.Split(r.name, ".")
		var codeDir string
		if len(parts) <= 1 {
			codeDir = srcDir
		} else {
			codeDir = filepath.Join(srcDir, filepath.Join(parts[:len(parts)-1]...))
		}
		shortName := parts[len(parts)-1]
		if err := ScaffoldFiles(codeDir, shortName, ext); err != nil {
			t.Fatalf("ScaffoldFiles(%q): %v", r.name, err)
		}
	}

	// Leaves should exist
	assertFileExists(t, filepath.Join(srcDir, "hello_world", "greet.ts"))
	assertFileExists(t, filepath.Join(srcDir, "hello_world", "greet.test.ts"))
	assertFileExists(t, filepath.Join(srcDir, "std", "io", "write_stdout.ts"))
	assertFileExists(t, filepath.Join(srcDir, "std", "io", "write_stdout.test.ts"))

	// Packages should NOT have files
	assertFileNotExists(t, filepath.Join(srcDir, "hello_world.ts"))
	assertFileNotExists(t, filepath.Join(srcDir, "hello_world.test.ts"))
	assertFileNotExists(t, filepath.Join(srcDir, "std.ts"))
	assertFileNotExists(t, filepath.Join(srcDir, "std.test.ts"))
	assertFileNotExists(t, filepath.Join(srcDir, "std", "io.ts"))
	assertFileNotExists(t, filepath.Join(srcDir, "std", "io.test.ts"))
}

// isLeaf mirrors rune.IsLeaf for the codegen test package.
func isLeaf(name string, allNames []string) bool {
	prefix := name + "."
	for _, n := range allNames {
		if strings.HasPrefix(n, prefix) {
			return false
		}
	}
	return true
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file %s to exist: %v", filepath.Base(path), err)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected file %s to NOT exist, but it does", path)
	}
}
