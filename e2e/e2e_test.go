package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "valkyrie-e2e-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "valkyrie")
	out, err := exec.Command("go", "build", "-o", binaryPath, "..").CombinedOutput()
	if err != nil {
		panic("failed to build binary: " + string(out))
	}

	os.Exit(m.Run())
}

func run(t *testing.T, configContent string) (stdout string, exitCode int) {
	t.Helper()
	f, err := os.CreateTemp("", "valkyrie-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "VALKYRIE_CONFIG="+f.Name())
	out, err := cmd.CombinedOutput()
	stdout = strings.TrimSpace(string(out))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdout, exitErr.ExitCode()
		}
	}
	return stdout, 0
}

func TestLoadsAgents(t *testing.T) {
	out, code := run(t, `
[agents.claude-api]
type = "claude-api"
model = "claude-sonnet-4-5"
api_key_env = "ANTHROPIC_API_KEY"

[agents.claude-max]
type = "claude-pro"
model = "claude-sonnet-4-5"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "2 agent(s) configured") {
		t.Errorf("expected agent count in output, got: %s", out)
	}
}

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

func TestClaudeProAgentConfig(t *testing.T) {
	out, code := run(t, `
[agents.pro]
type = "claude-pro"
model = "claude-sonnet-4-5"
credentials_path = "/tmp/fake-creds.json"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "pro") {
		t.Errorf("expected agent name 'pro' in output, got: %s", out)
	}
	if !strings.Contains(out, "claude-pro") {
		t.Errorf("expected type 'claude-pro' in output, got: %s", out)
	}
}

func TestDockerAgent(t *testing.T) {
	out, code := run(t, `
[agents.local]
type = "docker"
image = "ubuntu:22.04"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "docker") {
		t.Errorf("expected type 'docker' in output, got: %s", out)
	}
}
