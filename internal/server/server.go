package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/app"
	"github.com/chrishayen/odek/internal/decomposer"
	"github.com/chrishayen/odek/internal/feature"
	"github.com/chrishayen/odek/internal/hydrator"
	runepkg "github.com/chrishayen/odek/internal/rune"
	"github.com/chrishayen/odek/internal/server/jobs"
)

// Server is the Odek HTTP API server.
type Server struct {
	cfg          *config.Config
	runeStore    *runepkg.Store
	featureStore *feature.Store
	appStore     *app.Store
	dec          *decomposer.Decomposer
	hyd          *hydrator.Hydrator
	jobs         *jobs.Manager
	mux          *http.ServeMux
}

// New creates a new Server.
func New(cfg *config.Config, runeStore *runepkg.Store, featureStore *feature.Store, appStore *app.Store, dec *decomposer.Decomposer, hyd *hydrator.Hydrator) *Server {
	s := &Server{
		cfg:          cfg,
		runeStore:    runeStore,
		featureStore: featureStore,
		appStore:     appStore,
		dec:          dec,
		hyd:          hyd,
		jobs:         &jobs.Manager{},
		mux:          http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
	s.mux.HandleFunc("GET /api/runes", s.handleRunesList)
	s.mux.HandleFunc("GET /api/runes/{path...}", s.handleRunesGet)
	s.mux.HandleFunc("POST /api/decompose", s.handleDecompose)
	s.mux.HandleFunc("GET /api/decompose/{id}", s.handleJobStatus)
	s.mux.HandleFunc("POST /api/hydrate", s.handleHydrate)
	s.mux.HandleFunc("GET /api/hydrate/{id}", s.handleJobStatus)
	s.mux.HandleFunc("POST /api/ask", s.handleAsk)
	s.mux.HandleFunc("GET /api/ask/{id}", s.handleJobStatus)
	s.mux.HandleFunc("POST /api/check", s.handleCheck)
	s.mux.HandleFunc("POST /api/verify", s.handleVerify)
	s.mux.HandleFunc("GET /api/verify/{id}", s.handleJobStatus)
	s.mux.HandleFunc("GET /api/features", s.handleFeaturesList)
	s.mux.HandleFunc("GET /api/features/{name}", s.handleFeaturesGet)
	s.mux.HandleFunc("GET /api/apps", s.handleAppsList)
	s.mux.HandleFunc("GET /api/apps/{name}", s.handleAppsGet)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// progressWriter updates a Job's Progress field on each Write.
type progressWriter struct {
	jobs  *jobs.Manager
	jobID string
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.jobs.SetProgress(pw.jobID, strings.TrimSpace(string(p)))
	return len(p), nil
}

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleRunesList(w http.ResponseWriter, r *http.Request) {
	runes, err := s.runeStore.List()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, runes)
}

func (s *Server) handleRunesGet(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	dotPath := strings.ReplaceAll(path, "/", ".")
	rn, err := s.runeStore.Get(dotPath)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, rn)
}

func (s *Server) handleDecompose(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Requirement  string `json:"requirement"`
		Decomposition string `json:"decomposition,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Requirement == "" {
		jsonError(w, http.StatusBadRequest, "requirement is required")
		return
	}

	j := s.jobs.Create()
	go func() {
		s.jobs.SetRunning(j.ID)
		pw := &progressWriter{jobs: s.jobs, jobID: j.ID}
		result, err := s.dec.Decompose(context.Background(), req.Requirement, req.Decomposition, pw)
		if err != nil {
			s.jobs.SetFailed(j.ID, err)
			return
		}
		data, _ := json.Marshal(result)
		s.jobs.SetCompleted(j.ID, data)
	}()

	jsonResponse(w, http.StatusAccepted, map[string]string{"job_id": j.ID})
}

func (s *Server) handleAsk(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question string `json:"question"`
		Context  string `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Question == "" {
		jsonError(w, http.StatusBadRequest, "question is required")
		return
	}

	j := s.jobs.Create()
	go func() {
		s.jobs.SetRunning(j.ID)
		answer, err := s.dec.Ask(context.Background(), req.Question, req.Context)
		if err != nil {
			s.jobs.SetFailed(j.ID, err)
			return
		}
		data, _ := json.Marshal(answer)
		s.jobs.SetCompleted(j.ID, data)
	}()

	jsonResponse(w, http.StatusAccepted, map[string]string{"job_id": j.ID})
}

func (s *Server) handleHydrate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Verify      bool `json:"verify"`
		Concurrency int  `json:"concurrency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Concurrency == 0 {
		req.Concurrency = s.cfg.Concurrency
	}

	j := s.jobs.Create()
	go func() {
		s.jobs.SetRunning(j.ID)
		result, err := s.hyd.HydrateAll(context.Background(), req.Concurrency, req.Verify, nil)
		if err != nil {
			s.jobs.SetFailed(j.ID, err)
			return
		}
		data, _ := json.Marshal(result)
		s.jobs.SetCompleted(j.ID, data)
	}()

	jsonResponse(w, http.StatusAccepted, map[string]string{"job_id": j.ID})
}

func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	stale, ok, err := s.runeStore.CheckStaleRefs()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"stale": stale, "ok": ok})
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Concurrency int `json:"concurrency"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Concurrency == 0 {
		req.Concurrency = s.cfg.Concurrency
	}

	j := s.jobs.Create()
	go func() {
		s.jobs.SetRunning(j.ID)
		result, err := s.hyd.VerifyAll(context.Background(), req.Concurrency, nil)
		if err != nil {
			s.jobs.SetFailed(j.ID, err)
			return
		}
		data, _ := json.Marshal(result)
		s.jobs.SetCompleted(j.ID, data)
	}()

	jsonResponse(w, http.StatusAccepted, map[string]string{"job_id": j.ID})
}

func (s *Server) handleJobStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, ok := s.jobs.Get(id)
	if !ok {
		jsonError(w, http.StatusNotFound, "job not found")
		return
	}
	jsonResponse(w, http.StatusOK, j)
}

func (s *Server) handleFeaturesList(w http.ResponseWriter, r *http.Request) {
	features, err := s.featureStore.List()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, features)
}

func (s *Server) handleFeaturesGet(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	f, err := s.featureStore.Get(name)
	if err != nil {
		jsonError(w, http.StatusNotFound, fmt.Sprintf("feature %q not found", name))
		return
	}
	jsonResponse(w, http.StatusOK, f)
}

func (s *Server) handleAppsList(w http.ResponseWriter, r *http.Request) {
	apps, err := s.appStore.List()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, apps)
}

func (s *Server) handleAppsGet(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	a, err := s.appStore.Get(name)
	if err != nil {
		jsonError(w, http.StatusNotFound, fmt.Sprintf("app %q not found", name))
		return
	}
	jsonResponse(w, http.StatusOK, a)
}
