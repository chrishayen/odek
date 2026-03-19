package e2e_test

import (
	"net/http"
	"testing"
)

func TestClaudeAPIAgentStartsServer(t *testing.T) {
	base, cleanup := startServer(t, `
[agents.my-api]
type = "claude-api"
model = "claude-sonnet-4-5"
api_key_env = "ANTHROPIC_API_KEY"
`)
	defer cleanup()

	resp, err := http.Get(base + "/health")
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("server did not start with claude-api config")
	}
}

func TestClaudeProAgentStartsServer(t *testing.T) {
	base, cleanup := startServer(t, `
[agents.pro]
type = "claude-pro"
model = "claude-sonnet-4-5"
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	resp, err := http.Get(base + "/health")
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("server did not start with claude-pro config")
	}
}

func TestDockerAgentStartsServer(t *testing.T) {
	base, cleanup := startServer(t, `
[agents.local]
type = "docker"
image = "ubuntu:22.04"
`)
	defer cleanup()

	resp, err := http.Get(base + "/health")
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("server did not start with docker config")
	}
}
