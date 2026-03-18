package e2e_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMissingConfigEnv(t *testing.T) {
	cmd := exec.Command(binaryPath)
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
	out, code := run(t, `this is not valid toml ][[[`)
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid TOML\noutput: %s", out)
	}
}

func TestEmptyConfig(t *testing.T) {
	out, code := run(t, ``)
	if code != 0 {
		t.Fatalf("expected exit 0 for empty config, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "0 agent(s) configured") {
		t.Errorf("expected 0 agents, got: %s", out)
	}
}

func TestUnknownAgentType(t *testing.T) {
	out, code := run(t, `
[agents.bad]
type = "not-a-real-type"
`)
	if code == 0 {
		t.Fatalf("expected non-zero exit for unknown agent type\noutput: %s", out)
	}
	if !strings.Contains(out, "unknown type") {
		t.Errorf("expected 'unknown type' in error, got: %s", out)
	}
}

func TestLoadsMultipleAgents(t *testing.T) {
	out, code := run(t, `
[agents.claude-api]
type = "claude-api"
model = "claude-sonnet-4-5"
api_key_env = "ANTHROPIC_API_KEY"

[agents.claude-pro]
type = "claude-pro"
model = "claude-sonnet-4-5"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "2 agent(s) configured") {
		t.Errorf("expected 2 agents in output, got: %s", out)
	}
}
