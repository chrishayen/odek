package e2e_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMissingConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "valkyrie-empty-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Point at an empty dir — no config.toml exists
	out, code := run(t, tmp, "runes", "list")
	if code == 0 {
		t.Fatal("expected non-zero exit when config is missing")
	}
	if !strings.Contains(out, "config not found") {
		t.Errorf("expected 'config not found' in error, got: %s", out)
	}
}

func TestMissingConfigDirEnv(t *testing.T) {
	// Unset VALKYRIE_CONFIG_DIR — should fall back to ~/.config/valkyrie
	cmd := exec.Command(binaryPath, "runes", "list")
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "VALKYRIE_CONFIG_DIR=") {
			env = append(env, e)
		}
	}
	// Override HOME to a temp dir so it doesn't find a real config
	tmp, _ := os.MkdirTemp("", "valkyrie-nohome-*")
	defer os.RemoveAll(tmp)
	env = append(env, "HOME="+tmp)
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when no config exists")
	}
	if !strings.Contains(string(out), "config") {
		t.Errorf("expected config-related error, got: %s", string(out))
	}
}

func TestInvalidTOML(t *testing.T) {
	tmp, err := os.MkdirTemp("", "valkyrie-badtoml-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	os.WriteFile(tmp+"/config.toml", []byte("this is not valid toml ][[["), 0644)

	out, code := run(t, tmp, "runes", "list")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid TOML\noutput: %s", out)
	}
}

func TestUnknownAgentType(t *testing.T) {
	tmp, err := os.MkdirTemp("", "valkyrie-badagent-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	cfg := "[auth]\ndisabled = true\n\n[agents.bad]\ntype = \"not-a-real-type\"\n"
	os.WriteFile(tmp+"/config.toml", []byte(cfg), 0644)

	out, code := run(t, tmp, "runes", "list")
	if code == 0 {
		t.Fatalf("expected non-zero exit for unknown agent type\noutput: %s", out)
	}
	if !strings.Contains(out, "unknown type") {
		t.Errorf("expected 'unknown type' in error, got: %s", out)
	}
}
