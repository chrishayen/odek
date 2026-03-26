package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chrishayen/valkyrie/internal/adaptor"
	"github.com/chrishayen/valkyrie/internal/server/jobs"
	"github.com/chrishayen/valkyrie/internal/server/store"
)

type Server struct {
	store    *store.RuneStore
	token    string
	mux      *http.ServeMux
	jobs     *jobs.Manager
	pipeline *jobs.Pipeline
}

func New(dataDir, token string, a adaptor.Adaptor) *Server {
	st := store.NewRuneStore(dataDir)
	s := &Server{
		store:    st,
		token:    token,
		mux:      http.NewServeMux(),
		jobs:     jobs.NewManager(),
		pipeline: jobs.NewPipeline(a, st),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	log.Printf("valkyrie server listening on %s", addr)
	return http.ListenAndServe(addr, s)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)

	// Runes
	s.mux.HandleFunc("GET /api/runes", s.auth(s.handleRunesList))
	s.mux.HandleFunc("POST /api/runes", s.auth(s.handleRunesCreate))
	s.mux.HandleFunc("GET /api/runes/{fqn...}", s.auth(s.handleRunesGet))
	s.mux.HandleFunc("PUT /api/runes/{fqn...}", s.auth(s.handleRunesUpdate))
	s.mux.HandleFunc("DELETE /api/runes/{fqn...}", s.auth(s.handleRunesDelete))
	s.mux.HandleFunc("POST /api/runes_search", s.auth(s.handleRunesSearch))
	s.mux.HandleFunc("POST /api/runes_commit/{fqn...}", s.auth(s.handleRunesCommit))

	// Projects
	s.mux.HandleFunc("GET /api/projects", s.auth(s.handleProjectsList))
	s.mux.HandleFunc("GET /api/projects/{name}", s.auth(s.handleProjectRunes))

	// Requirements
	s.mux.HandleFunc("POST /api/requirements", s.auth(s.handleRequirementsSubmit))
	s.mux.HandleFunc("GET /api/requirements/{id}", s.auth(s.handleRequirementsStatus))
}

// auth wraps a handler with bearer token authentication.
func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.token != "" {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != s.token {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- Rune handlers ---

func (s *Server) handleRunesList(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	namespace := r.URL.Query().Get("namespace")
	runes, err := s.store.List(project, namespace)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if runes == nil {
		runes = []store.Rune{}
	}
	writeJSON(w, http.StatusOK, runes)
}

func (s *Server) handleRunesCreate(w http.ResponseWriter, r *http.Request) {
	var rune store.Rune
	if err := json.NewDecoder(r.Body).Decode(&rune); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if err := s.store.Create(rune); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := s.store.Get(rune.FQN)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleRunesGet(w http.ResponseWriter, r *http.Request) {
	fqn := fqnFromPath(r.PathValue("fqn"))
	rune, err := s.store.Get(fqn)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rune)
}

func (s *Server) handleRunesUpdate(w http.ResponseWriter, r *http.Request) {
	fqn := fqnFromPath(r.PathValue("fqn"))
	existing, err := s.store.Get(fqn)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var patch struct {
		Description *string  `json:"description,omitempty"`
		Signature   *string  `json:"signature,omitempty"`
		Behavior    *string  `json:"behavior,omitempty"`
		Version     *string  `json:"version,omitempty"`
		Status      *string  `json:"status,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if patch.Description != nil {
		existing.Description = *patch.Description
	}
	if patch.Signature != nil {
		existing.Signature = *patch.Signature
	}
	if patch.Behavior != nil {
		existing.Behavior = *patch.Behavior
	}
	if patch.Version != nil {
		existing.Version = *patch.Version
	}
	if patch.Status != nil {
		existing.Status = *patch.Status
	}

	if err := s.store.Update(*existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updated, _ := s.store.Get(fqn)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleRunesDelete(w http.ResponseWriter, r *http.Request) {
	fqn := fqnFromPath(r.PathValue("fqn"))
	if err := s.store.Delete(fqn); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": fqn})
}

func (s *Server) handleRunesSearch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	runes, err := s.store.Search(body.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if runes == nil {
		runes = []store.Rune{}
	}
	writeJSON(w, http.StatusOK, runes)
}

func (s *Server) handleRunesCommit(w http.ResponseWriter, r *http.Request) {
	fqn := fqnFromPath(r.PathValue("fqn"))
	existing, err := s.store.Get(fqn)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	existing.Status = "approved"
	if err := s.store.Update(*existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, existing)
}

// --- Project handlers ---

func (s *Server) handleProjectsList(w http.ResponseWriter, _ *http.Request) {
	projects, err := s.store.ListProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if projects == nil {
		projects = []string{}
	}
	writeJSON(w, http.StatusOK, projects)
}

func (s *Server) handleProjectRunes(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	runes, err := s.store.List("", name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if runes == nil {
		runes = []store.Rune{}
	}
	writeJSON(w, http.StatusOK, runes)
}

// --- Requirements handlers ---

func (s *Server) handleRequirementsSubmit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Project      string `json:"project"`
		Requirements string `json:"requirements"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if body.Requirements == "" {
		writeError(w, http.StatusBadRequest, "requirements is required")
		return
	}

	job := s.jobs.Create()
	go func() {
		s.jobs.SetRunning(job.ID)
		result, err := s.pipeline.Run(r.Context(), body.Project, body.Requirements)
		if err != nil {
			s.jobs.SetFailed(job.ID, err)
			return
		}
		data, _ := json.Marshal(result)
		s.jobs.SetCompleted(job.ID, data)
	}()

	writeJSON(w, http.StatusAccepted, map[string]string{
		"id":     job.ID,
		"status": string(job.Status),
	})
}

func (s *Server) handleRequirementsStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job := s.jobs.Get(id)
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

// --- Helpers ---

// fqnFromPath converts URL path segments "net/http/parse_url" to dot notation "net.http.parse_url".
func fqnFromPath(path string) string {
	return strings.ReplaceAll(path, "/", ".")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Store returns the underlying rune store for direct access by other server components.
func (s *Server) Store() *store.RuneStore {
	return s.store
}

// Addr returns a formatted address string.
func Addr(port int) string {
	return fmt.Sprintf(":%d", port)
}
