package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/app"
	"github.com/chrishayen/odek/internal/claude"
	"github.com/chrishayen/odek/internal/decomposer"
	"github.com/chrishayen/odek/internal/hydrator"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

func newTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	dir := t.TempDir()
	outPath := filepath.Join(dir, "src")
	cfg := &config.Config{
		Project:      "test",
		Language:     "go",
		RegistryPath: dir,
		OutputPath:   outPath,
		Concurrency:  2,
		Agent:        config.Agent{Mock: true},
		Server:       config.Server{Port: 0},
	}
	rs := runepkg.NewStore(dir, outPath)
	as := app.NewStore(dir, outPath)
	client := claude.New("", "", true)
	dec := decomposer.New(rs, client)
	hyd := hydrator.New(rs, client, "go")
	s := New(cfg, rs, as, dec, hyd)
	return s, dir
}

func TestHealthEndpoint(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("status = %q", resp["status"])
	}
}

func TestRunesListEmpty(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/runes", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
}

func TestRunesGetNotFound(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/runes/no/such/rune", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestRunesGetFound(t *testing.T) {
	s, dir := newTestServer(t)
	rs := runepkg.NewStore(dir, filepath.Join(dir, "src"))
	rs.Create(runepkg.Rune{
		Name:        "test.hello",
		Description: "says hello",
		Signature:   "(name: string) -> string",
	})

	req := httptest.NewRequest("GET", "/api/runes/test/hello", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
}

func TestFeaturesListEmpty(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/features", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
}

func TestFeaturesGetNotFound(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/features/nope", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestAppsListEmpty(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/apps", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
}

func TestAppsGetNotFound(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/apps/nope", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestCheckEndpoint(t *testing.T) {
	s, _ := newTestServer(t)
	req := httptest.NewRequest("POST", "/api/check", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
}
