package e2e_test

import (
	"testing"
)

func TestClaudeSubAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agent]
type = "claude-sub"
model = "claude-sonnet-4-5"
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with claude-sub agent should load successfully")
	}
}

func TestMockAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agent]
type = "mock"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with mock agent should load successfully")
	}
}

func TestDefaultAgentType(t *testing.T) {
	// No [agent] section — should default to claude-sub
	dir, cleanup := testEnv(t, `
[agent]
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config without explicit agent type should default to claude-sub")
	}
}
