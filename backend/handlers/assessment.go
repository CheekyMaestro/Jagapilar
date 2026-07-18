package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"jagapilar-backend/database"
	"jagapilar-backend/middleware"
	"jagapilar-backend/models"
	"jagapilar-backend/services"
)

// GetAssessmentItems handles GET /api/assessment/items?pillar=parent
func GetAssessmentItems(w http.ResponseWriter, r *http.Request) {
	pillar := r.URL.Query().Get("pillar")

	var query string
	var args []interface{}

	if pillar != "" {
		query = "SELECT code, pillar, question_text, is_reverse, construct, reference_source FROM assessment_items WHERE pillar = $1 ORDER BY code"
		args = append(args, pillar)
	} else {
		query = "SELECT code, pillar, question_text, is_reverse, construct, reference_source FROM assessment_items ORDER BY code"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil item asesmen")
		return
	}
	defer rows.Close()

	var items []models.AssessmentItem
	for rows.Next() {
		var item models.AssessmentItem
		var refSource sql.NullString
		if err := rows.Scan(&item.Code, &item.Pillar, &item.QuestionText, &item.IsReverse, &item.Construct, &refSource); err != nil {
			continue
		}
		if refSource.Valid {
			item.ReferenceSource = refSource.String
		}
		items = append(items, item)
	}

	if items == nil {
		items = []models.AssessmentItem{}
	}

	respondJSON(w, http.StatusOK, items)
}

// CreateSession handles POST /api/assessment/sessions
func CreateSession(w http.ResponseWriter, r *http.Request) {
	neuroID := r.Header.Get("X-Neuro-ID")
	pillar := r.Header.Get("X-User-Role") // Should match the intended pillar
	
	// Get auth context if present
	userID, okUser := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.RoleKey).(string)

	if neuroID == "" {
		respondError(w, http.StatusBadRequest, "Neuro ID diperlukan")
		return
	}

	validPillars := map[string]bool{"parent": true, "teacher": true, "student": true}
	if !validPillars[pillar] {
		respondError(w, http.StatusBadRequest, "Pillar harus parent, teacher, atau student")
		return
	}

	// ENFORCEMENT: Parent can only start parent session, Teacher -> teacher
	if pillar != "student" {
		if !okUser {
			respondError(w, http.StatusUnauthorized, "Sesi kadaluarsa, silakan login ulang")
			return
		}
		if userRole != pillar {
			respondError(w, http.StatusForbidden, "Akses ditolak: Peran Anda tidak sesuai dengan asesmen ini")
			return
		}
	}

	// Resolve child_id from neuro_id
	var childID string
	err := database.DB.QueryRow("SELECT id FROM children WHERE neuro_id = $1", neuroID).Scan(&childID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Neuro ID tidak valid")
		return
	}

	var session models.AssessmentSession
	
	// Insert session
	var userIDPtr *string
	if okUser {
		userIDPtr = &userID
	}
	
	err = database.DB.QueryRow(
		`INSERT INTO assessment_sessions (child_id, user_id, pillar, status)
		 VALUES ($1, $2, $3, 'draft')
		 RETURNING id, child_id, user_id, pillar, status, created_at`,
		childID, userIDPtr, pillar,
	).Scan(&session.ID, &session.ChildID, &session.UserID, &session.Pillar, &session.Status, &session.CreatedAt)

	if err != nil {
		log.Printf("Error creating session: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal membuat sesi asesmen")
		return
	}

	respondJSON(w, http.StatusCreated, session)
}

// SaveResponse handles POST /api/assessment/responses (progressive save)
func SaveResponse(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RawScore < 1 || req.RawScore > 5 {
		respondError(w, http.StatusBadRequest, "Skor harus antara 1-5")
		return
	}

	// Upsert response
	_, err := database.DB.Exec(
		`INSERT INTO item_responses (session_id, item_code, raw_score)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (session_id, item_code) DO UPDATE SET raw_score = $3, saved_at = now()`,
		req.SessionID, req.ItemCode, req.RawScore,
	)
	if err != nil {
		// Fallback if unique constraint doesn't exist
		_, err = database.DB.Exec(
			`INSERT INTO item_responses (session_id, item_code, raw_score)
			 VALUES ($1, $2, $3)`,
			req.SessionID, req.ItemCode, req.RawScore,
		)
		if err != nil {
			log.Printf("Error saving response: %v", err)
			respondError(w, http.StatusInternalServerError, "Gagal menyimpan jawaban")
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

// SubmitAssessment handles POST /api/assessment/submit
func SubmitAssessment(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// JWT Role validation
	_, okUser := r.Context().Value(middleware.UserIDKey).(string)
	userRole, _ := r.Context().Value(middleware.RoleKey).(string)

	if req.Pillar != "student" {
		if !okUser {
			respondError(w, http.StatusUnauthorized, "Sesi kadaluarsa, silakan login ulang")
			return
		}
		if userRole != req.Pillar {
			respondError(w, http.StatusForbidden, "Akses ditolak: Anda tidak diizinkan submit asesmen pilar ini")
			return
		}
	}

	// Resolve child_id from neuro_id
	var childID string
	err := database.DB.QueryRow("SELECT id FROM children WHERE neuro_id = $1", req.NeuroID).Scan(&childID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Neuro ID tidak valid")
		return
	}

	if len(req.Items) > 0 {
		var responses []models.ItemResponse
		for _, item := range req.Items {
			responses = append(responses, models.ItemResponse{
				ItemCode: item.ItemCode,
				RawScore: item.RawScore,
			})
		}

		// Save pillar score
		pillarScore, err := services.CalculateAndSavePillarScore(childID, req.Pillar, responses)
		if err != nil {
			log.Printf("Error calculating pillar score: %v", err)
			respondError(w, http.StatusInternalServerError, "Gagal menghitung skor")
			return
		}

		// Check if all 3 pillars have been submitted
		var pillarCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM pillar_scores WHERE child_id = $1", childID).Scan(&pillarCount)

		var compositeResult *models.CompositeScore
		if pillarCount >= 3 {
			compositeResult, err = services.OnAllThreePillarsSubmitted(childID)
			if err != nil {
				log.Printf("Error calculating composite: %v", err)
			}
		}

		result := map[string]interface{}{
			"pillar_score":  pillarScore,
			"pillars_done":  pillarCount,
			"all_complete":  pillarCount >= 3,
		}
		if compositeResult != nil {
			result["composite"] = compositeResult
		}

		respondJSON(w, http.StatusOK, result)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "received",
		"pillar":      req.Pillar,
		"total_score": req.TotalScore,
		"zone":        req.Zone,
	})
}

// GetResults handles GET /api/assessment/results/{neuro_id}
func GetResults(w http.ResponseWriter, r *http.Request) {
	neuroID := extractPathParam(r.URL.Path, "/api/assessment/results/")
	if neuroID == "" {
		respondError(w, http.StatusBadRequest, "Neuro ID diperlukan")
		return
	}

	var childID string
	err := database.DB.QueryRow("SELECT id FROM children WHERE neuro_id = $1", neuroID).Scan(&childID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Neuro ID tidak ditemukan")
		return
	}

	result, err := services.GetAssessmentResult(childID)
	if err != nil {
		log.Printf("Error fetching results: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil hasil asesmen")
		return
	}

	logAudit("viewer", "read", "assessment_results", childID)
	respondJSON(w, http.StatusOK, result)
}

// AssessmentHandler routes assessment requests
func AssessmentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/assessment/items" && r.Method == http.MethodGet:
		GetAssessmentItems(w, r)

	case path == "/api/assessment/sessions" && r.Method == http.MethodPost:
		CreateSession(w, r)

	case path == "/api/assessment/responses" && r.Method == http.MethodPost:
		SaveResponse(w, r)

	case path == "/api/assessment/submit" && r.Method == http.MethodPost:
		SubmitAssessment(w, r)

	case strings.HasPrefix(path, "/api/assessment/results/") && r.Method == http.MethodGet:
		GetResults(w, r)

	default:
		respondError(w, http.StatusNotFound, "Endpoint not found")
	}
}
