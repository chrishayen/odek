package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/codegen"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

func TestCommitScaffoldsOnlyLeaves(t *testing.T) {
	dir := t.TempDir()
	runesDir := filepath.Join(dir, "runes")
	srcDir := filepath.Join(dir, "src")
	store := runepkg.NewStore(runesDir, srcDir)
	suffix := "a1b2c3"
	lang := "ts"

	// Create draft runes (same as saveDraft does)
	drafts := []runepkg.Rune{
		{Name: "hello_world_" + suffix, Description: "feature root", Status: "draft"},
		{Name: "hello_world_" + suffix + ".greet", Description: "greets", Signature: "(name: string) -> string", Status: "draft"},
		{Name: "std_" + suffix, Description: "standard lib", Status: "draft"},
		{Name: "std_" + suffix + ".io", Description: "io package", Status: "draft"},
		{Name: "std_" + suffix + ".io.write_stdout", Description: "writes stdout", Signature: "(msg: string) -> void", Status: "draft"},
	}
	for _, r := range drafts {
		if err := store.Create(r); err != nil {
			t.Fatalf("Create draft %q: %v", r.Name, err)
		}
	}

	// ---- Reproduce the commit path exactly ----
	allDrafts, err := store.ListByStatus("draft")
	if err != nil {
		t.Fatal(err)
	}

	// Collect all names (store + batch) for leaf detection
	allRunes, _ := store.List()
	allNames := make([]string, 0, len(allRunes)+len(allDrafts))
	for _, r := range allRunes {
		allNames = append(allNames, r.Name)
	}
	for _, r := range allDrafts {
		if hasDraftSuffix(r.Name, suffix) {
			allNames = append(allNames, removeDraftSuffix(r.Name, suffix))
		}
	}

	t.Logf("allNames: %v", allNames)

	ext := config.LangExtension(lang)

	for _, r := range allDrafts {
		if !hasDraftSuffix(r.Name, suffix) {
			continue
		}
		clean := r
		clean.Name = removeDraftSuffix(r.Name, suffix)
		clean.Status = ""
		_ = store.Delete(clean.Name)
		if err := store.Create(clean); err != nil {
			t.Fatalf("Create clean %q: %v", clean.Name, err)
		}

		isLeaf := runepkg.IsLeaf(clean.Name, allNames)
		t.Logf("  %q isLeaf=%v codeDir=%q shortName=%q", clean.Name, isLeaf, store.CodeDir(clean.Name), runepkg.ShortName(clean.Name))

		if isLeaf {
			codegen.ScaffoldFiles(store.CodeDir(clean.Name), runepkg.ShortName(clean.Name), ext)
		}
		_ = store.Delete(r.Name)
	}

	// Assert leaf files exist
	assertExists(t, filepath.Join(srcDir, "hello_world", "greet.ts"))
	assertExists(t, filepath.Join(srcDir, "hello_world", "greet.test.ts"))
	assertExists(t, filepath.Join(srcDir, "std", "io", "write_stdout.ts"))
	assertExists(t, filepath.Join(srcDir, "std", "io", "write_stdout.test.ts"))

	// Assert package files do NOT exist
	assertNotExists(t, filepath.Join(srcDir, "hello_world.ts"))
	assertNotExists(t, filepath.Join(srcDir, "hello_world.test.ts"))
	assertNotExists(t, filepath.Join(srcDir, "std.ts"))
	assertNotExists(t, filepath.Join(srcDir, "std.test.ts"))
	assertNotExists(t, filepath.Join(srcDir, "std", "io.ts"))
	assertNotExists(t, filepath.Join(srcDir, "std", "io.test.ts"))
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to exist: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected %s to NOT exist, but it does", path)
	}
}
