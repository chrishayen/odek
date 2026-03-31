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
	tmp, err := os.MkdirTemp("", "odek-e2e-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "odek")
	out, err := exec.Command("go", "build", "-o", binaryPath, "..").CombinedOutput()
	if err != nil {
		panic("failed to build binary: " + string(out))
	}

	os.Exit(m.Run())
}

// testEnv creates an isolated project dir with odek.toml and returns the dir
// path and a cleanup function. The registry lives inside the project dir.
func testEnv(t *testing.T, extraTOML string) (projectDir string, cleanup func()) {
	t.Helper()

	projectDir, err := os.MkdirTemp("", "odek-env-*")
	if err != nil {
		t.Fatal(err)
	}

	agentTOML := extraTOML
	if agentTOML == "" {
		agentTOML = "[agent]\nmock = true\n"
	}
	toml := "project = \"test-project\"\n\n" + agentTOML
	if err := os.WriteFile(filepath.Join(projectDir, "odek.toml"), []byte(toml), 0644); err != nil {
		t.Fatal(err)
	}

	cleanup = func() { os.RemoveAll(projectDir) }
	return projectDir, cleanup
}

// run executes the odek binary with the given project dir and args.
// Returns combined output and exit code.
func run(t *testing.T, projectDir string, args ...string) (output string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "VALKYRIE_PROJECT_DIR="+projectDir)
	out, err := cmd.CombinedOutput()
	output = strings.TrimSpace(string(out))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, exitErr.ExitCode()
		}
	}
	return output, 0
}
