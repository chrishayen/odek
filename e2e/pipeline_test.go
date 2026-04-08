package e2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// configYAMLPath returns the absolute path to config.yaml in the project root.
func configYAMLPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "config.yaml")
}

// startProxy launches "odek serve" in the background against the given project
// dir and waits for the proxy to become reachable on port 8317. Returns a
// cleanup function that kills the process.
func startProxy(t *testing.T, projectDir string) func() {
	t.Helper()

	cmd := exec.Command(binaryPath, "serve")
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(),
		"ODEK_PROJECT_DIR="+projectDir,
		"CLIPROXY_CONFIG="+configYAMLPath(),
	)
	logFile, _ := os.CreateTemp("", "proxy-*.log")
	defer os.Remove(logFile.Name())
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start proxy: %v", err)
	}

	// Poll until proxy is reachable (needs API key header).
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8317/v1/models", nil)
		req.Header.Set("x-api-key", "sk-local-proxy")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return func() { cmd.Process.Kill(); cmd.Wait() }
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	cmd.Process.Kill()
	cmd.Wait()
	logFile.Seek(0, 0)
	logBytes, _ := os.ReadFile(logFile.Name())
	t.Fatalf("proxy did not become ready within 20s\nconfig: %s\nlog:\n%s", configYAMLPath(), string(logBytes))
	return nil
}

func TestFullPipelineDecomposeToExport(t *testing.T) {
	dir, cleanup := testEnv(t, "language = \"ts\"\n")
	defer cleanup()

	stopProxy := startProxy(t, dir)
	defer stopProxy()

	// Step 1: Decompose — LLM determines the rune structure.
	out, code := run(t, dir, "runes", "decompose", "--yes",
		"Say hello world to the user")
	if code != 0 {
		t.Fatalf("decompose failed (exit %d): %s", code, out)
	}
	if !strings.Contains(out, "created rune") {
		t.Fatalf("expected at least one rune to be created, got: %s", out)
	}

	// Step 2: Hydrate all runes.
	out, code = run(t, dir, "runes", "hydrate-all")
	if code != 0 {
		t.Fatalf("hydrate-all failed (exit %d): %s", code, out)
	}
	if strings.Contains(out, `"failed": 0`) == false {
		t.Errorf("expected no hydration failures, got: %s", out)
	}

	// Step 3: Discover the feature name from the runes list (top-level rune
	// that isn't "std") and compose it.
	out, code = run(t, dir, "runes", "list")
	if code != 0 {
		t.Fatalf("runes list failed (exit %d): %s", code, out)
	}
	var runes []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(out), &runes); err != nil {
		t.Fatalf("failed to parse runes list: %v\nraw: %s", err, out)
	}
	var featureName string
	for _, r := range runes {
		if !strings.Contains(r.Name, ".") && r.Name != "std" {
			featureName = r.Name
			break
		}
	}
	if featureName == "" {
		t.Fatal("could not find top-level feature rune after decompose")
	}

	out, code = run(t, dir, "features", "compose", featureName)
	if code != 0 {
		t.Fatalf("compose failed (exit %d): %s", code, out)
	}

	// Step 4: Export.
	out, code = run(t, dir, "features", "export", featureName)
	if code != 0 {
		t.Fatalf("export failed (exit %d): %s", code, out)
	}

	var exportResult struct {
		FeatureName string   `json:"feature_name"`
		Version     string   `json:"version"`
		OutputDir   string   `json:"output_dir"`
		Files       []string `json:"files"`
	}
	if err := json.Unmarshal([]byte(out), &exportResult); err != nil {
		t.Fatalf("failed to parse export output: %v\nraw: %s", err, out)
	}
	if exportResult.FeatureName == "" {
		t.Error("expected feature_name in export result")
	}
	if len(exportResult.Files) == 0 {
		t.Error("expected non-empty files list in export result")
	}

	// Verify dist structure.
	distDir := filepath.Join(dir, "dist", featureName)
	for _, f := range []string{"package.json", "index.ts"} {
		if _, err := os.Stat(filepath.Join(distDir, f)); err != nil {
			t.Errorf("expected %s in export: %v", f, err)
		}
	}

	// Verify at least one .ts source file was exported.
	hasTS := false
	for _, f := range exportResult.Files {
		if strings.HasSuffix(f, ".ts") && f != "index.ts" {
			hasTS = true
			break
		}
	}
	if !hasTS {
		t.Errorf("expected at least one .ts source file in export, got: %v", exportResult.Files)
	}

	fmt.Printf("Feature %q exported %d files: %v\n", featureName, len(exportResult.Files), exportResult.Files)
}
