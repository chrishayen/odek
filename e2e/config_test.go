package e2e_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMissingConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "odek-empty-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Point at an empty dir — no odek.toml exists
	out, code := run(t, tmp, "runes", "list")
	if code == 0 {
		t.Fatal("expected non-zero exit when config is missing")
	}
	if !strings.Contains(out, "odek.toml not found") {
		t.Errorf("expected 'odek.toml not found' in error, got: %s", out)
	}
}

func TestMissingConfigNoEnv(t *testing.T) {
	// Unset ODEK_PROJECT_DIR and run from an empty temp dir
	tmp, err := os.MkdirTemp("", "odek-nohome-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	cmd := exec.Command(binaryPath, "runes", "list")
	cmd.Dir = tmp
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "ODEK_PROJECT_DIR=") {
			env = append(env, e)
		}
	}
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when no config exists")
	}
	if !strings.Contains(string(out), "odek.toml not found") {
		t.Errorf("expected odek.toml error, got: %s", string(out))
	}
}

func TestInvalidTOML(t *testing.T) {
	tmp, err := os.MkdirTemp("", "odek-badtoml-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	os.WriteFile(tmp+"/odek.toml", []byte("this is not valid toml ][[["), 0644)

	out, code := run(t, tmp, "runes", "list")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid TOML\noutput: %s", out)
	}
}

