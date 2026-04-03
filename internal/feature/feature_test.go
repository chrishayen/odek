package feature

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreCreateGetRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	if err := s.Create("auth", "Authentication feature"); err != nil {
		t.Fatal(err)
	}

	f, err := s.Get("auth")
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "auth" {
		t.Errorf("Name = %q", f.Name)
	}
	if f.Version != "0.1.0" {
		t.Errorf("Version = %q", f.Version)
	}
	if f.Status != "draft" {
		t.Errorf("Status = %q", f.Status)
	}
}

func TestStoreCreateDuplicate(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("auth", "first")
	if err := s.Create("auth", "second"); err == nil {
		t.Error("expected error for duplicate")
	}
}

func TestStoreCreateValidation(t *testing.T) {
	s := NewStore(t.TempDir(), "")

	if err := s.Create("", "body"); err == nil {
		t.Error("expected error for empty name")
	}
	if err := s.Create("has space", "body"); err == nil {
		t.Error("expected error for name with space")
	}
	if err := s.Create("auth", ""); err == nil {
		t.Error("expected error for empty body")
	}
}

func TestStoreList(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("auth", "Authentication")
	s.Create("payment", "Payment processing")

	features, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 2 {
		t.Errorf("List() len = %d", len(features))
	}
}

func TestStoreListEmpty(t *testing.T) {
	s := NewStore(t.TempDir(), "")
	features, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if features != nil {
		t.Errorf("expected nil for empty registry, got %v", features)
	}
}

func TestStoreUpdate(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("auth", "Authentication")
	f, _ := s.Get("auth")
	f.Version = "1.0.0"
	f.Status = "stable"

	if err := s.Update(*f); err != nil {
		t.Fatal(err)
	}

	updated, _ := s.Get("auth")
	if updated.Version != "1.0.0" {
		t.Errorf("Version = %q", updated.Version)
	}
	if updated.Status != "stable" {
		t.Errorf("Status = %q", updated.Status)
	}
}

func TestStoreDelete(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("auth", "Authentication")
	if err := s.Delete("auth"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("auth"); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestStoreDeleteNotFound(t *testing.T) {
	s := NewStore(t.TempDir(), "")
	if err := s.Delete("nope"); err == nil {
		t.Error("expected error for missing feature")
	}
}

func TestStoreGetNotFound(t *testing.T) {
	s := NewStore(t.TempDir(), "")
	_, err := s.Get("nope")
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestReadRaw(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("auth", "Auth feature body")
	raw, err := s.ReadRaw("auth")
	if err != nil {
		t.Fatal(err)
	}
	if raw == "" {
		t.Error("expected non-empty raw content")
	}
}

func TestCodeDir(t *testing.T) {
	s := NewStore("/registry", "/out")
	got := s.CodeDir("auth")
	want := filepath.Join("/out", "auth")
	if got != want {
		t.Errorf("CodeDir() = %q, want %q", got, want)
	}
}

func TestListSkipsOutputDir(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "src")
	s := NewStore(dir, outDir)

	s.Create("auth", "Auth feature")

	// Create a directory matching the output path with a feature.md
	os.MkdirAll(outDir, 0755)
	os.WriteFile(filepath.Join(outDir, "feature.md"), []byte("---\nversion: 0.1.0\nstatus: draft\n---\n\nfake"), 0644)

	features, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 1 {
		t.Errorf("expected 1 feature (output dir excluded), got %d", len(features))
	}
}
