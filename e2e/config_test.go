package e2e_test

import (
	"os"
	"strings"
	"testing"
)

func TestMissingConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "valkyrie-empty-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	out, code := run(t, tmp, "mcp")
	if code == 0 {
		t.Fatal("expected non-zero exit when config is missing")
	}
	if !strings.Contains(out, "valkyrie.toml not found") {
		t.Errorf("expected 'valkyrie.toml not found' in error, got: %s", out)
	}
}

func TestInvalidTOML(t *testing.T) {
	tmp, err := os.MkdirTemp("", "valkyrie-badtoml-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	os.WriteFile(tmp+"/valkyrie.toml", []byte("this is not valid toml ][[["), 0644)

	out, code := run(t, tmp, "mcp")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid TOML\noutput: %s", out)
	}
}
