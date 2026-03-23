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

// testEnv creates an isolated config dir with config.toml and returns the dir
// path and a cleanup function. The registry lives inside the config dir.
func testEnv(t *testing.T, extraTOML string) (configDir string, cleanup func()) {
	t.Helper()

	configDir, err := os.MkdirTemp("", "valkyrie-env-*")
	if err != nil {
		t.Fatal(err)
	}

	registryDir := filepath.Join(configDir, "registry")
	if err := os.MkdirAll(registryDir, 0755); err != nil {
		t.Fatal(err)
	}

	toml := "registry_path = " + quote(registryDir) + "\n\n[auth]\ndisabled = true\n\n" + extraTOML
	if err := os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(toml), 0644); err != nil {
		t.Fatal(err)
	}

	cleanup = func() { os.RemoveAll(configDir) }
	return configDir, cleanup
}

// run executes the valkyrie binary with the given config dir and args.
// Returns combined output and exit code.
func run(t *testing.T, configDir string, args ...string) (output string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "VALKYRIE_CONFIG_DIR="+configDir)
	out, err := cmd.CombinedOutput()
	output = strings.TrimSpace(string(out))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, exitErr.ExitCode()
		}
	}
	return output, 0
}

func quote(s string) string {
	return `"` + s + `"`
}
