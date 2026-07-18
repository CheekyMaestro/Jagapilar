package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"jagapilar-backend/database"
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
	childID := r.Header.Get("X-Child-ID")
	informantID := r.Header.Get("X-Informant-ID")
	pillar := r.Header.Get("X-User-Role")

	if childID == "" || informantID == "" {
		respondError(w, http.StatusUnauthorized, "Sesi tidak valid, gunakan link akses yang benar")
		return
	}

	validPillars := map[string]bool{"parent": true, "teacher": true, "student": true}
	if !validPillars[pillar] {
		respondError(w, http.StatusBadRequest, "Pillar harus parent, teacher, atau student")
		return
	}

	var session models.AssessmentSession
	err := database.DB.QueryRow(
		`INSERT INTO assessment_sessions (child_id, informant_id, pillar, status)
		 VALUES ($1, $2, $3, 'draft')
		 RETURNING id, child_id, informant_id, pillar, status, created_at`,
		childID, informantID, pillar,
	).Scan(&session.ID, &session.ChildID, &session.InformantID, &session.Pillar, &session.Status, &session.CreatedAt)

	if err != nil {
		log.Printf("Error creating session: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal membuat sesi asesmen")
		return
	}

	log.Printf("📋 Assessment session created: %s (pillar=%s)", session.ID, session.Pillar)
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

	// Upsert response (replace if already answered)
	_, err := database.DB.Exec(
		`INSERT INTO item_responses (session_id, item_code, raw_score)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (session_id, item_code) DO UPDATE SET raw_score = $3, saved_at = now()`,
		req.SessionID, req.ItemCode, req.RawScore,
	)
	// If unique constraint doesn't exist for session_id + item_code, use simple insert
	if err != nil {
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
// This finalizes a session and triggers scoring
func SubmitAssessment(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Override with header values from token if present
	if hChildID := r.Header.Get("X-Child-ID"); hChildID != "" {
		req.ChildID = hChildID
	}
	if hPillar := r.Header.Get("X-User-Role"); hPillar != "" {
		req.Pillar = hPillar
	}

	// If child_id is provided, save the scores server-side
	if req.ChildID != "" && len(req.Items) > 0 {
		// Convert items to ItemResponse format
		var responses []models.ItemResponse
		for _, item := range req.Items {
			responses = append(responses, models.ItemResponse{
				ItemCode: item.ItemCode,
				RawScore: item.RawScore,
			})
		}

		// Calculate and save pillar score
		pillarScore, err := services.CalculateAndSavePillarScore(req.ChildID, req.Pillar, responses)
		if err != nil {
			log.Printf("Error calculating pillar score: %v", err)
			respondError(w, http.StatusInternalServerError, "Gagal menghitung skor")
			return
		}

		// Check if all 3 pillars have been submitted
		var pillarCount int
		database.DB.QueryRow(
			"SELECT COUNT(*) FROM pillar_scores WHERE child_id = $1",
			req.ChildID,
		).Scan(&pillarCount)

		var compositeResult *models.CompositeScore
		if pillarCount >= 3 {
			compositeResult, err = services.OnAllThreePillarsSubmitted(req.ChildID)
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

	// Simple acknowledgment if no child_id (client-side scoring mode)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "received",
		"pillar":      req.Pillar,
		"total_score": req.TotalScore,
		"zone":        req.Zone,
	})
}

// GetResults handles GET /api/assessment/results/{child_id}
func GetResults(w http.ResponseWriter, r *http.Request) {
	childID := extractPathParam(r.URL.Path, "/api/assessment/results/")
	if childID == "" {
		respondError(w, http.StatusBadRequest, "Child ID required")
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
