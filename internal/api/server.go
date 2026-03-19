package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	runepkg "github.com/chrishayen/valkyrie/internal/rune"
)

type Server struct {
	store *runepkg.Store
	mux   *http.ServeMux
}

func NewServer(store *runepkg.Store) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/runes", s.runesCollection)
	s.mux.HandleFunc("/runes/", s.runesItem)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) runesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		runes, err := s.store.List()
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		if runes == nil {
			runes = []runepkg.Rune{}
		}
		writeJSON(w, http.StatusOK, runes)

	case http.MethodPost:
		var rune runepkg.Rune
		if err := readJSON(r, &rune); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		if err := s.store.Create(rune); err != nil {
			writeError(w, err, clientOrServer(err))
			return
		}
		created, _ := s.store.Get(rune.Name)
		writeJSON(w, http.StatusCreated, created)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) runesItem(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// /runes/{name}
	if len(parts) == 2 {
		name := parts[1]
		switch r.Method {
		case http.MethodGet:
			rune, err := s.store.Get(name)
			if err != nil {
				writeError(w, err, http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, rune)

		case http.MethodPut:
			var rune runepkg.Rune
			if err := readJSON(r, &rune); err != nil {
				writeError(w, err, http.StatusBadRequest)
				return
			}
			rune.Name = name
			if err := s.store.Update(rune); err != nil {
				writeError(w, err, clientOrServer(err))
				return
			}
			updated, _ := s.store.Get(name)
			writeJSON(w, http.StatusOK, updated)

		case http.MethodDelete:
			if err := s.store.Delete(name); err != nil {
				writeError(w, err, http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func readJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("request body required")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func clientOrServer(err error) int {
	msg := err.Error()
	if strings.Contains(msg, "not found") ||
		strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "required") ||
		strings.Contains(msg, "slug") ||
		strings.Contains(msg, "already stable") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
