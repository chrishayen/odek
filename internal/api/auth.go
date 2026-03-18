package api

import (
	"net/http"
	"strings"
)

// authMiddleware enforces bearer token authentication unless disabled.
func authMiddleware(token string, disabled bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if disabled {
			next.ServeHTTP(w, r)
			return
		}

		// /health is always public
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] != token {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized — provide a valid Bearer token",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
