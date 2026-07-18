package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"jagapilar-backend/database"
	"jagapilar-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("JAGAPILAR_SECRET_KEY_CHANGE_ME_IN_PROD")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterUser handles POST /api/auth/register
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Role != "parent" && req.Role != "teacher" {
		respondError(w, http.StatusBadRequest, "Role harus parent atau teacher")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal memproses password")
		return
	}

	var user models.User
	err = database.DB.QueryRow(
		`INSERT INTO users (name, email_contact, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email_contact, role, created_at`,
		req.Name, req.EmailContact, string(hashedPassword), req.Role,
	).Scan(&user.ID, &user.Name, &user.EmailContact, &user.Role, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			respondError(w, http.StatusConflict, "Email/Kontak sudah terdaftar")
			return
		}
		log.Printf("Register error: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mendaftarkan akun")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// LoginUser handles POST /api/auth/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	var passwordHash string
	var schoolID sql.NullString

	err := database.DB.QueryRow(
		`SELECT id, name, email_contact, password_hash, role, school_id, created_at
		 FROM users WHERE email_contact = $1`,
		req.EmailContact,
	).Scan(&user.ID, &user.Name, &user.EmailContact, &passwordHash, &user.Role, &schoolID, &user.CreatedAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusUnauthorized, "Email atau Password salah")
		return
	} else if err != nil {
		log.Printf("Login error: %v", err)
		respondError(w, http.StatusInternalServerError, "Terjadi kesalahan sistem")
		return
	}

	if schoolID.Valid {
		user.SchoolID = &schoolID.String
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "Email atau Password salah")
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal membuat sesi")
		return
	}

	respondJSON(w, http.StatusOK, models.AuthResponse{
		Token: tokenString,
		User:  user,
	})
}

// AuthHandler routes auth requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/auth/register":
		RegisterUser(w, r)
	case path == "/api/auth/login":
		LoginUser(w, r)
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
