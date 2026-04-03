package app

import (
	"path/filepath"
	"testing"
)

func TestStoreCreateGetRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	if err := s.Create("myapp", "A test application"); err != nil {
		t.Fatal(err)
	}

	a, err := s.Get("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "myapp" {
		t.Errorf("Name = %q", a.Name)
	}
	if a.Version != "0.1.0" {
		t.Errorf("Version = %q", a.Version)
	}
	if a.Status != "draft" {
		t.Errorf("Status = %q", a.Status)
	}
}

func TestStoreCreateDuplicate(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("myapp", "first")
	if err := s.Create("myapp", "second"); err == nil {
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
	if err := s.Create("myapp", ""); err == nil {
		t.Error("expected error for empty body")
	}
}

func TestStoreList(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("app1", "first app")
	s.Create("app2", "second app")

	apps, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Errorf("List() len = %d", len(apps))
	}
}

func TestStoreUpdate(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("myapp", "A test app")
	a, _ := s.Get("myapp")
	a.Version = "1.0.0"
	a.Status = "stable"
	a.EntryPoint = "auth"

	if err := s.Update(*a); err != nil {
		t.Fatal(err)
	}

	updated, _ := s.Get("myapp")
	if updated.Version != "1.0.0" {
		t.Errorf("Version = %q", updated.Version)
	}
	if updated.Status != "stable" {
		t.Errorf("Status = %q", updated.Status)
	}
	if updated.EntryPoint != "auth" {
		t.Errorf("EntryPoint = %q", updated.EntryPoint)
	}
}

func TestStoreDelete(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	s.Create("myapp", "A test app")
	if err := s.Delete("myapp"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("myapp"); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestStoreGetNotFound(t *testing.T) {
	s := NewStore(t.TempDir(), "")
	_, err := s.Get("nope")
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestCodeDir(t *testing.T) {
	s := NewStore("/registry", "/out")
	got := s.CodeDir("myapp")
	want := filepath.Join("/out", "myapp")
	if got != want {
		t.Errorf("CodeDir() = %q, want %q", got, want)
	}
}
