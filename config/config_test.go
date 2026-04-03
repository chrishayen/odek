package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLangExtension(t *testing.T) {
	tests := []struct {
		lang, want string
	}{
		{"go", ".go"},
		{"py", ".py"},
		{"ts", ".ts"},
		{"unknown", ".ts"},
	}
	for _, tt := range tests {
		if got := LangExtension(tt.lang); got != tt.want {
			t.Errorf("LangExtension(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

func TestLangName(t *testing.T) {
	tests := []struct {
		lang, want string
	}{
		{"go", "Go"},
		{"py", "Python"},
		{"ts", "TypeScript (Node.js with node:* built-ins)"},
		{"unknown", "TypeScript (Node.js with node:* built-ins)"},
	}
	for _, tt := range tests {
		if got := LangName(tt.lang); got != tt.want {
			t.Errorf("LangName(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	toml := `project = "mytest"
language = "go"

[agent]
mock = true
`
	os.WriteFile(filepath.Join(dir, "odek.toml"), []byte(toml), 0644)
	t.Setenv("ODEK_PROJECT_DIR", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Language != "go" {
		t.Errorf("Language = %q", cfg.Language)
	}
	if cfg.Concurrency != 50 {
		t.Errorf("Concurrency = %d", cfg.Concurrency)
	}
	if cfg.Server.Port != 8319 {
		t.Errorf("Server.Port = %d", cfg.Server.Port)
	}
	if !cfg.Agent.Mock {
		t.Error("expected Agent.Mock = true")
	}
}

func TestLoadUnsupportedLanguage(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "odek.toml"), []byte(`project = "x"
language = "ruby"
`), 0644)
	t.Setenv("ODEK_PROJECT_DIR", dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for unsupported language")
	}
}

func TestLoadDefaultLanguage(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "odek.toml"), []byte(`project = "x"`), 0644)
	t.Setenv("ODEK_PROJECT_DIR", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Language != "go" {
		t.Errorf("expected default language go, got %q", cfg.Language)
	}
}

func TestFindRootWalksUp(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "odek.toml"), []byte(`project = "test"`), 0644)

	subdir := filepath.Join(root, "a", "b", "c")
	os.MkdirAll(subdir, 0755)
	t.Setenv("ODEK_PROJECT_DIR", subdir)

	found, err := FindRoot()
	if err != nil {
		t.Fatal(err)
	}
	if found != root {
		t.Errorf("FindRoot() = %q, want %q", found, root)
	}
}

func TestFindRootNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ODEK_PROJECT_DIR", dir)

	_, err := FindRoot()
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestResolveTokenPrecedence(t *testing.T) {
	a := Agent{Token: "literal"}
	if got := a.ResolveToken(); got != "literal" {
		t.Errorf("literal token: got %q", got)
	}

	a = Agent{TokenEnv: "TEST_ODEK_TOKEN"}
	t.Setenv("TEST_ODEK_TOKEN", "from-env")
	if got := a.ResolveToken(); got != "from-env" {
		t.Errorf("env token: got %q", got)
	}

	a = Agent{}
	t.Setenv("CLAUDE_CODE_OAUTH_TOKEN", "oauth-tok")
	if got := a.ResolveToken(); got != "oauth-tok" {
		t.Errorf("oauth token: got %q", got)
	}
}
