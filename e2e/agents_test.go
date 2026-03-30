package e2e_test

import (
	"testing"
)

func TestAgentConfigWithToken(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agent]
model = "claude-sonnet-4-5"
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with token should load successfully")
	}
}

func TestMockAgentConfigLoads(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agent]
mock = true
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with mock agent should load successfully")
	}
}

func TestDefaultAgentConfig(t *testing.T) {
	dir, cleanup := testEnv(t, `
[agent]
token = "sk-ant-oat01-fake"
`)
	defer cleanup()

	_, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Error("config with default agent settings should load successfully")
	}
}
