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

// run executes the binary with the given TOML config content.
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
