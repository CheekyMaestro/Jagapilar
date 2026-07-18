package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"jagapilar-backend/database"
	"jagapilar-backend/models"
)

// CreateChild handles POST /api/children
func CreateChild(w http.ResponseWriter, r *http.Request) {
	var req models.CreateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SchoolID == "" {
		respondError(w, http.StatusBadRequest, "school_id wajib diisi")
		return
	}

	// Generate anonymous code (JAGAPILAR-XXXXX)
	anonCode := generateAnonCode()

	var child models.Child
	err := database.DB.QueryRow(
		`INSERT INTO children (anon_code, school_id, grade, birth_year)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, anon_code, school_id, grade, birth_year, created_at`,
		anonCode, req.SchoolID, req.Grade, req.BirthYear,
	).Scan(&child.ID, &child.AnonCode, &child.SchoolID, &child.Grade, &child.BirthYear, &child.CreatedAt)

	if err != nil {
		log.Printf("Error creating child: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mendaftarkan anak")
		return
	}

	logAudit("admin", "create", "children", child.ID)
	log.Printf("👶 Child registered: %s (grade %d)", child.AnonCode, child.Grade)
	respondJSON(w, http.StatusCreated, child)
}

// GetChild handles GET /api/children/{id}
func GetChild(w http.ResponseWriter, r *http.Request) {
	childID := extractPathParam(r.URL.Path, "/api/children/")
	if childID == "" {
		respondError(w, http.StatusBadRequest, "Child ID required")
		return
	}

	// Also remove sub-paths like /informants
	if strings.Contains(childID, "/") {
		childID = strings.Split(childID, "/")[0]
	}

	var child models.Child
	err := database.DB.QueryRow(
		`SELECT id, anon_code, school_id, grade, birth_year, created_at
		 FROM children WHERE id = $1`,
		childID,
	).Scan(&child.ID, &child.AnonCode, &child.SchoolID, &child.Grade, &child.BirthYear, &child.CreatedAt)

	if err != nil {
		respondError(w, http.StatusNotFound, "Anak tidak ditemukan")
		return
	}

	respondJSON(w, http.StatusOK, child)
}

// CreateInformant handles POST /api/children/{id}/informants
func CreateInformant(w http.ResponseWriter, r *http.Request) {
	// Extract child ID from path
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/children/"), "/")
	if len(parts) < 2 {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	childID := parts[0]

	var req models.CreateInformantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate role
	validRoles := map[string]bool{"parent": true, "teacher": true, "student": true}
	if !validRoles[req.Role] {
		respondError(w, http.StatusBadRequest, "Role harus parent, teacher, atau student")
		return
	}

	// Generate unique access token (magic link style)
	accessToken := generateAccessToken()
	tokenExpires := time.Now().Add(7 * 24 * time.Hour) // 7 days

	var informant models.Informant
	err := database.DB.QueryRow(
		`INSERT INTO informants (child_id, role, contact_hash, access_token, token_expires_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, child_id, role, contact_hash, access_token, token_expires_at, created_at`,
		childID, req.Role, req.ContactHash, accessToken, tokenExpires,
	).Scan(&informant.ID, &informant.ChildID, &informant.Role, &informant.ContactHash,
		&informant.AccessToken, &informant.TokenExpiresAt, &informant.CreatedAt)

	if err != nil {
		log.Printf("Error creating informant: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal membuat informan")
		return
	}

	logAudit(req.Role, "create", "informants", informant.ID)
	log.Printf("👤 Informant created: role=%s child=%s token=%s", req.Role, childID, accessToken[:8]+"...")
	respondJSON(w, http.StatusCreated, informant)
}

// RegisterChildFull handles POST /api/children/register-full
func RegisterChildFull(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterChildFullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SchoolID == "" {
		respondError(w, http.StatusBadRequest, "school_id wajib diisi")
		return
	}

	anonCode := generateAnonCode()

	// Use transaction
	tx, err := database.DB.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}
	defer tx.Rollback()

	var child models.Child
	err = tx.QueryRow(
		`INSERT INTO children (anon_code, school_id, grade, birth_year)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, anon_code, school_id, grade, birth_year, created_at`,
		anonCode, req.SchoolID, req.Grade, req.BirthYear,
	).Scan(&child.ID, &child.AnonCode, &child.SchoolID, &child.Grade, &child.BirthYear, &child.CreatedAt)

	if err != nil {
		log.Printf("Error creating child: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mendaftarkan anak")
		return
	}

	roles := []string{"parent", "teacher", "student"}
	var tokens [3]string

	tokenExpires := time.Now().Add(7 * 24 * time.Hour) // 7 days

	for i, role := range roles {
		accessToken := generateAccessToken()
		tokens[i] = accessToken
		
		_, err = tx.Exec(
			`INSERT INTO informants (child_id, role, access_token, token_expires_at)
			 VALUES ($1, $2, $3, $4)`,
			child.ID, role, accessToken, tokenExpires,
		)
		if err != nil {
			log.Printf("Error creating informant %s: %v", role, err)
			respondError(w, http.StatusInternalServerError, "Gagal membuat informan")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menyimpan data")
		return
	}

	logAudit("admin", "register_full", "children", child.ID)
	log.Printf("👶 Child registered fully: %s (grade %d) with 3 tokens", child.AnonCode, child.Grade)

	respondJSON(w, http.StatusCreated, models.ChildWithTokens{
		Child:        child,
		ParentToken:  tokens[0],
		TeacherToken: tokens[1],
		StudentToken: tokens[2],
	})
}

// ChildrenHandler routes children requests
func ChildrenHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// POST /api/children/register-full
	if path == "/api/children/register-full" && r.Method == http.MethodPost {
		RegisterChildFull(w, r)
		return
	}

	// POST /api/children/{id}/informants
	if strings.Contains(path, "/informants") && r.Method == http.MethodPost {
		CreateInformant(w, r)
		return
	}

	// GET /api/children/{id}
	if path != "/api/children" && path != "/api/children/" {
		if r.Method == http.MethodGet {
			GetChild(w, r)
			return
		}
	}

	// POST /api/children
	if r.Method == http.MethodPost {
		CreateChild(w, r)
		return
	}

	respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

// generateAnonCode creates a unique anonymous identifier
func generateAnonCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("JP-%s", strings.ToUpper(hex.EncodeToString(b)))
}

// generateAccessToken creates a secure random access token
func generateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
