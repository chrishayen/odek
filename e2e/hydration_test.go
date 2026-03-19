package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

const helloWorldDescription = "Returns the string 'Hello, World!' when called"

func startServerWithMock(t *testing.T) (baseURL string, registryDir string, cleanup func()) {
	t.Helper()
	return startServerFull(t, `
[agents.mock]
type = "mock"
`, true, "")
}

func hydrateRune(t *testing.T, base, runeName, sandbox string) map[string]any {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"sandbox": sandbox})
	resp, err := http.Post(base+"/runes/"+runeName+"/hydrate", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("hydrate request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		t.Fatalf("hydrate returned %d: %v", resp.StatusCode, errBody)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode hydrate result: %v", err)
	}
	return result
}

func createRune(t *testing.T, base, name string) {
	t.Helper()
	resp := apiPost(t, base+"/runes", map[string]any{
		"name":        name,
		"description": helloWorldDescription,
	})
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		t.Fatalf("create rune %q: expected 201, got %d: %v", name, resp.StatusCode, errBody)
	}
}

func TestHydrateHelloWorld(t *testing.T) {
	base, _, cleanup := startServerWithMock(t)
	defer cleanup()

	createRune(t, base, "hello-world-single")
	result := hydrateRune(t, base, "hello-world-single", "mock")

	if result["RuneName"] != "hello-world-single" && result["rune_name"] != "hello-world-single" {
		t.Errorf("expected rune name in hydrate result, got %v", result)
	}
	if _, ok := result["Coverage"]; !ok {
		if _, ok2 := result["coverage"]; !ok2 {
			t.Errorf("expected coverage in hydrate result, got %v", result)
		}
	}

	resp := apiGet(t, base+"/runes/hello-world-single")
	defer resp.Body.Close()
	var rune map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&rune)
	if rune["hydrated"] != true {
		t.Errorf("expected hydrated=true, got %v", rune["hydrated"])
	}
}

func TestHydrateFourSessions(t *testing.T) {
	type session struct {
		base    string
		cleanup func()
	}

	sessions := make([]session, 4)
	for i := range sessions {
		base, _, cleanup := startServerWithMock(t)
		sessions[i] = session{base: base, cleanup: cleanup}
	}
	defer func() {
		for _, s := range sessions {
			s.cleanup()
		}
	}()

	for i, s := range sessions {
		createRune(t, s.base, fmt.Sprintf("hello-world-%d", i+1))
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	for i, s := range sessions {
		wg.Add(1)
		go func(idx int, base string) {
			defer wg.Done()
			result := hydrateRune(t, base, fmt.Sprintf("hello-world-%d", idx+1), "mock")
			if result["RuneName"] == nil && result["rune_name"] == nil {
				errCh <- fmt.Errorf("session %d: missing rune name in result", idx)
			}
		}(i, s.base)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
}

func TestHydrateGeneratesCodeFiles(t *testing.T) {
	base, registryDir, cleanup := startServerWithMock(t)
	defer cleanup()

	createRune(t, base, "hello-world-files")
	hydrateRune(t, base, "hello-world-files", "mock")

	codeDir := filepath.Join(registryDir, "runes", "hello-world-files")
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		t.Fatalf("code dir not created: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated files")
	}
}

func TestHydrateCoverageTracked(t *testing.T) {
	base, _, cleanup := startServerWithMock(t)
	defer cleanup()

	createRune(t, base, "hello-world-coverage")
	hydrateRune(t, base, "hello-world-coverage", "mock")

	resp := apiGet(t, base+"/runes/hello-world-coverage")
	defer resp.Body.Close()
	var rune map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&rune)
	if _, ok := rune["coverage"]; !ok {
		t.Errorf("expected coverage field on rune, got %v", rune)
	}
}
