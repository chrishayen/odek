package feature

import (
	"path/filepath"
	"testing"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

func TestCodeDir(t *testing.T) {
	rs := runepkg.NewStore("/registry", "/out")
	s := NewStore(rs, "/out")
	got := s.CodeDir("auth")
	want := filepath.Join("/out", "auth")
	if got != want {
		t.Errorf("CodeDir() = %q, want %q", got, want)
	}
}

func TestListDerivesFromRunes(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(dir, filepath.Join(dir, "src"))
	s := NewStore(rs, filepath.Join(dir, "src"))

	// Create a top-level rune (feature) and a nested rune
	rs.Create(runepkg.Rune{
		Name:        "auth",
		Description: "Authentication",
		Signature:   "() -> void",
	})
	rs.Create(runepkg.Rune{
		Name:        "auth.login",
		Description: "Login handler",
		Signature:   "(user: string) -> bool",
	})

	features, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 1 {
		t.Errorf("List() len = %d, want 1 (only top-level)", len(features))
	}
	if len(features) > 0 && features[0].Name != "auth" {
		t.Errorf("Name = %q, want %q", features[0].Name, "auth")
	}
}

func TestListEmpty(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(dir, filepath.Join(dir, "src"))
	s := NewStore(rs, filepath.Join(dir, "src"))

	features, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 0 {
		t.Errorf("expected 0 features, got %d", len(features))
	}
}

func TestGetNotFound(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(dir, filepath.Join(dir, "src"))
	s := NewStore(rs, filepath.Join(dir, "src"))

	_, err := s.Get("nope")
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestGetDerivesFromRune(t *testing.T) {
	dir := t.TempDir()
	rs := runepkg.NewStore(dir, filepath.Join(dir, "src"))
	s := NewStore(rs, filepath.Join(dir, "src"))

	rs.Create(runepkg.Rune{
		Name:        "auth",
		Description: "Authentication",
		Signature:   "() -> void",
	})

	f, err := s.Get("auth")
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "auth" {
		t.Errorf("Name = %q", f.Name)
	}
	if f.Version != "1.0.0" {
		t.Errorf("Version = %q", f.Version)
	}
}
