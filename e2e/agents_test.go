package e2e_test

import (
	"testing"
)

func TestClaudeAPIAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agents.my-api]
type = "claude-api"
model = "claude-sonnet-4-5"
api_key_env = "ANTHROPIC_API_KEY"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with claude-api agent should load successfully")
	}
}

func TestClaudeMaxAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agents.max]
type = "claude-max"
model = "claude-sonnet-4-5"
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with claude-max agent should load successfully")
	}
}

func TestDockerAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agents.local]
type = "docker"
image = "ubuntu:22.04"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with docker agent should load successfully")
	}
}
