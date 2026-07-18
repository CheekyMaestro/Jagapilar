package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"jagapilar-backend/database"
	"jagapilar-backend/models"
)

// ValidateToken handles POST /api/auth/validate-token
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		respondJSON(w, http.StatusOK, models.TokenValidationResponse{Valid: false})
		return
	}

	var informant models.Informant
	err := database.DB.QueryRow(
		`SELECT id, child_id, role, access_token, token_expires_at
		 FROM informants WHERE access_token = $1`,
		req.Token,
	).Scan(&informant.ID, &informant.ChildID, &informant.Role, &informant.AccessToken, &informant.TokenExpiresAt)

	if err != nil {
		respondJSON(w, http.StatusOK, models.TokenValidationResponse{Valid: false})
		return
	}

	// Check expiration
	if informant.TokenExpiresAt != nil && informant.TokenExpiresAt.Before(time.Now()) {
		respondJSON(w, http.StatusOK, models.TokenValidationResponse{Valid: false})
		return
	}

	expiresStr := ""
	if informant.TokenExpiresAt != nil {
		expiresStr = informant.TokenExpiresAt.Format(time.RFC3339)
	}

	logAudit(informant.Role, "token_validated", "informants", informant.ID)

	respondJSON(w, http.StatusOK, models.TokenValidationResponse{
		Valid:     true,
		Role:      informant.Role,
		ChildID:   informant.ChildID,
		ExpiresAt: expiresStr,
	})
}

// AuthHandler routes auth requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/auth/validate-token":
		ValidateToken(w, r)
	default:
		respondError(w, http.StatusNotFound, "Endpoint not found")
	}
}

// ==================== Shared Helpers ====================

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error JSON response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.APIResponse{
		Success: false,
		Error:   message,
	})
}

// extractPathParam extracts a parameter from a URL path
func extractPathParam(path, prefix string) string {
	param := strings.TrimPrefix(path, prefix)
	param = strings.TrimSuffix(param, "/")
	// Remove sub-paths
	if idx := strings.Index(param, "/"); idx != -1 {
		param = param[:idx]
	}
	return param
}

// logAudit records an action in the audit log
func logAudit(actorRole, action, targetTable, targetID string) {
	_, err := database.DB.Exec(
		`INSERT INTO audit_log (actor_role, action, target_table, target_id)
		 VALUES ($1, $2, $3, $4)`,
		actorRole, action, targetTable, targetID,
	)
	if err != nil {
		log.Printf("⚠️ Audit log failed: %v", err)
	}
}
