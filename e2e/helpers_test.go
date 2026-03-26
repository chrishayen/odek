package e2e_test

import (
	"fmt"
	"net"
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

// testEnv creates an isolated project dir with valkyrie.toml and returns the dir
// path and a cleanup function.
func testEnv(t *testing.T, serverURL string) (projectDir string, cleanup func()) {
	t.Helper()

	projectDir, err := os.MkdirTemp("", "valkyrie-env-*")
	if err != nil {
		t.Fatal(err)
	}

	toml := fmt.Sprintf("project = \"test-project\"\n\n[server]\nurl = %q\ntoken_env = \"VALKYRIE_TOKEN\"\n", serverURL)
	if err := os.WriteFile(filepath.Join(projectDir, "valkyrie.toml"), []byte(toml), 0644); err != nil {
		t.Fatal(err)
	}

	cleanup = func() { os.RemoveAll(projectDir) }
	return projectDir, cleanup
}

// run executes the valkyrie binary with the given project dir and args.
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

// freePort finds a free TCP port.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

// startServer launches a valkyrie server on a random port and returns the URL and cleanup function.
func startServer(t *testing.T) (url string, cleanup func()) {
	t.Helper()
	port := freePort(t)
	dataDir, err := os.MkdirTemp("", "valkyrie-data-*")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "serve",
		"--port", fmt.Sprintf("%d", port),
		"--data-dir", dataDir,
		"--token", "test-token",
	)
	cmd.Env = append(os.Environ(), "ANTHROPIC_API_KEY=test-key-not-used")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	url = fmt.Sprintf("http://127.0.0.1:%d", port)

	// Wait for server to be ready
	for i := 0; i < 50; i++ {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	cleanup = func() {
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dataDir)
	}
	return url, cleanup
}
