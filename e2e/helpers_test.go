package e2e_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

var nextPort = 18200

func allocPort() int {
	nextPort++
	return nextPort
}

func writeTempTOML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "valkyrie-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

// startServer spins up `valkyrie serve` with auth disabled (for local/test use).
func startServer(t *testing.T, configTOML string) (baseURL string, cleanup func()) {
	t.Helper()
	return startServerRaw(t, configTOML, true, "")
}

// startServerWithToken spins up `valkyrie serve` with bearer token auth enabled.
func startServerWithToken(t *testing.T, token string) (baseURL string, cleanup func()) {
	t.Helper()
	return startServerRaw(t, "", false, token)
}

// startServerRaw is the underlying helper for starting the server with full control.
func startServerRaw(t *testing.T, configTOML string, authDisabled bool, authToken string) (baseURL string, cleanup func()) {
	t.Helper()

	registryDir, err := os.MkdirTemp("", "valkyrie-registry-*")
	if err != nil {
		t.Fatal(err)
	}

	var authSection string
	if authDisabled {
		authSection = "[auth]\ndisabled = true\n"
	} else {
		authSection = fmt.Sprintf("[auth]\ntoken = %q\n", authToken)
	}

	fullConfig := fmt.Sprintf("registry_path = %q\n\n%s\n%s", registryDir, authSection, configTOML)
	cfgFile := writeTempTOML(t, fullConfig)

	port := allocPort()
	addr := fmt.Sprintf(":%d", port)
	baseURL = fmt.Sprintf("http://localhost:%d", port)

	cmd := exec.Command(binaryPath, "serve", "--addr", addr)
	cmd.Env = append(os.Environ(), "VALKYRIE_CONFIG="+cfgFile)

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Wait for server to be ready
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return baseURL, func() {
		cmd.Process.Kill()
		os.RemoveAll(registryDir)
	}
}

// runBinary runs the binary with a config file and returns stdout+stderr and exit code.
func runBinary(t *testing.T, configContent string, args ...string) (output string, exitCode int) {
	t.Helper()
	cfgFile := writeTempTOML(t, configContent)
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "VALKYRIE_CONFIG="+cfgFile)
	out, err := cmd.CombinedOutput()
	output = strings.TrimSpace(string(out))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, exitErr.ExitCode()
		}
	}
	return output, 0
}
