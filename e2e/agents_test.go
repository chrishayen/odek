package e2e_test

import (
	"strings"
	"testing"
)

func TestClaudeAPIAgent(t *testing.T) {
	out, code := run(t, `
[agents.my-api]
type = "claude-api"
model = "claude-sonnet-4-5"
api_key_env = "ANTHROPIC_API_KEY"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "my-api") {
		t.Errorf("expected agent name in output, got: %s", out)
	}
	if !strings.Contains(out, "claude-api") {
		t.Errorf("expected type in output, got: %s", out)
	}
}

func TestClaudeProAgent(t *testing.T) {
	out, code := run(t, `
[agents.pro]
type = "claude-pro"
model = "claude-sonnet-4-5"
token = "sk-ant-oat01-fake"
`)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, out)
	}
	if !strings.Contains(out, "pro") {
		t.Errorf("expected agent name in output, got: %s", out)
	}
	if !strings.Contains(out, "claude-pro") {
		t.Errorf("expected type in output, got: %s", out)
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
	if !strings.Contains(out, "local") {
		t.Errorf("expected agent name in output, got: %s", out)
	}
	if !strings.Contains(out, "docker") {
		t.Errorf("expected type in output, got: %s", out)
	}
}
