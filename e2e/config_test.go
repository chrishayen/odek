package e2e_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMissingConfigEnv(t *testing.T) {
	cmd := exec.Command(binaryPath, "serve")
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "VALKYRIE_CONFIG=") {
			env = append(env, e)
		}
	}
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when VALKYRIE_CONFIG is unset")
	}
	if !strings.Contains(string(out), "VALKYRIE_CONFIG") {
		t.Errorf("expected error mentioning VALKYRIE_CONFIG, got: %s", string(out))
	}
}

func TestInvalidTOML(t *testing.T) {
	out, code := runBinary(t, `this is not valid toml ][[[`, "serve")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid TOML\noutput: %s", out)
	}
}

func TestUnknownAgentType(t *testing.T) {
	out, code := runBinary(t, `
[agents.bad]
type = "not-a-real-type"
`, "serve")
	if code == 0 {
		t.Fatalf("expected non-zero exit for unknown agent type\noutput: %s", out)
	}
	if !strings.Contains(out, "unknown type") {
		t.Errorf("expected 'unknown type' in error, got: %s", out)
	}
}
