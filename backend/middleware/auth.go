package middleware

import (
	"net/http"
	"strings"
	"time"

	"jagapilar-backend/database"
)

// Auth middleware validates token-based authentication
// Checks the Authorization header for a valid access token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		publicPaths := []string{
			"/api/schools",
			"/api/assessment/items",
			"/api/auth/",
		}
		for _, p := range publicPaths {
			if strings.HasPrefix(r.URL.Path, p) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Also skip for non-API paths (static files)
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Token diperlukan"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, `{"error":"Format token tidak valid"}`, http.StatusUnauthorized)
			return
		}

		// Validate token against database
		var informantID, childID, role string
		var expiresAt *time.Time
		err := database.DB.QueryRow(
			"SELECT id, child_id, role, token_expires_at FROM informants WHERE access_token = $1",
			token,
		).Scan(&informantID, &childID, &role, &expiresAt)

		if err != nil {
			http.Error(w, `{"error":"Token tidak valid"}`, http.StatusUnauthorized)
			return
		}

		// Check expiration
		if expiresAt != nil && expiresAt.Before(time.Now()) {
			http.Error(w, `{"error":"Token sudah kedaluwarsa"}`, http.StatusUnauthorized)
			return
		}

		// Set variables in header for downstream handlers
		r.Header.Set("X-Informant-ID", informantID)
		r.Header.Set("X-Child-ID", childID)
		r.Header.Set("X-User-Role", role)

		next.ServeHTTP(w, r)
	})
}
