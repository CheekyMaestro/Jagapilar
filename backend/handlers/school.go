package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"jagapilar-backend/database"
	"jagapilar-backend/models"
)

// CreateSchool handles POST /api/schools
func CreateSchool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.City == "" {
		respondError(w, http.StatusBadRequest, "Nama sekolah dan kota wajib diisi")
		return
	}

	var school models.School
	err := database.DB.QueryRow(
		`INSERT INTO schools (name, city, grade_level, principal_name, total_classes)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, city, grade_level, principal_name, total_classes, created_at`,
		req.Name, req.City, req.GradeLevel, req.PrincipalName, req.TotalClasses,
	).Scan(&school.ID, &school.Name, &school.City, &school.GradeLevel,
		&school.PrincipalName, &school.TotalClasses, &school.CreatedAt)

	if err != nil {
		log.Printf("Error creating school: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mendaftarkan sekolah")
		return
	}

	// Log to audit
	logAudit("admin", "create", "schools", school.ID)

	log.Printf("🏫 School registered: %s (%s)", school.Name, school.City)
	respondJSON(w, http.StatusCreated, school)
}

// ListSchools handles GET /api/schools
func ListSchools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, name, city, grade_level, principal_name, total_classes, created_at
		 FROM schools ORDER BY created_at DESC`,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data sekolah")
		return
	}
	defer rows.Close()

	var schools []models.School
	for rows.Next() {
		var s models.School
		if err := rows.Scan(&s.ID, &s.Name, &s.City, &s.GradeLevel,
			&s.PrincipalName, &s.TotalClasses, &s.CreatedAt); err != nil {
			continue
		}
		schools = append(schools, s)
	}

	if schools == nil {
		schools = []models.School{}
	}

	respondJSON(w, http.StatusOK, schools)
}

// GetSchoolDashboard handles GET /api/schools/{id}/dashboard
func GetSchoolDashboard(w http.ResponseWriter, r *http.Request) {
	schoolID := extractPathParam(r.URL.Path, "/api/schools/")
	
	// Remove /dashboard from the end
	if strings.HasSuffix(schoolID, "/dashboard") {
		schoolID = strings.TrimSuffix(schoolID, "/dashboard")
	}

	if schoolID == "" {
		respondError(w, http.StatusBadRequest, "School ID required")
		return
	}

	var resp models.SchoolDashboardResponse

	// Get School Info
	err := database.DB.QueryRow(
		`SELECT id, name, city, grade_level, principal_name, total_classes, created_at
		 FROM schools WHERE id = $1`,
		schoolID,
	).Scan(&resp.School.ID, &resp.School.Name, &resp.School.City, &resp.School.GradeLevel,
		&resp.School.PrincipalName, &resp.School.TotalClasses, &resp.School.CreatedAt)

	if err != nil {
		respondError(w, http.StatusNotFound, "Sekolah tidak ditemukan")
		return
	}

	// Get all children
	rows, err := database.DB.Query(
		`SELECT id, anon_code, grade, birth_year FROM children WHERE school_id = $1 ORDER BY created_at DESC`,
		schoolID,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data anak")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var childItem models.ChildDashboardItem
		childItem.Child.SchoolID = schoolID
		if err := rows.Scan(&childItem.Child.ID, &childItem.Child.AnonCode, &childItem.Child.Grade, &childItem.Child.BirthYear); err != nil {
			continue
		}

		// Get composite score for the child
		var comp models.CompositeScore
		err = database.DB.QueryRow(
			`SELECT composite, total_raw, composite_zone, needs_manual_review FROM composite_scores WHERE child_id = $1`,
			childItem.Child.ID,
		).Scan(&comp.Composite, &comp.TotalRaw, &comp.CompositeZone, &comp.NeedsManualReview)
		
		if err == nil {
			childItem.Composite = &comp
			// Update stats
			switch comp.CompositeZone {
			case "hijau":
				resp.TotalHijau++
			case "kuning":
				resp.TotalKuning++
			case "merah":
				resp.TotalMerah++
			}
			if comp.NeedsManualReview {
				resp.TotalReview++
			}
		}

		// Get tokens
		tokenRows, err := database.DB.Query(
			`SELECT role, access_token FROM informants WHERE child_id = $1`,
			childItem.Child.ID,
		)
		if err == nil {
			for tokenRows.Next() {
				var role, token string
				if err := tokenRows.Scan(&role, &token); err == nil {
					switch role {
					case "parent":
						childItem.ParentToken = token
					case "teacher":
						childItem.TeacherToken = token
					case "student":
						childItem.StudentToken = token
					}
				}
			}
			tokenRows.Close()
		}

		resp.Children = append(resp.Children, childItem)
	}

	if resp.Children == nil {
		resp.Children = []models.ChildDashboardItem{}
	}

	respondJSON(w, http.StatusOK, resp)
}

// SchoolsHandler routes school requests
func SchoolsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasSuffix(path, "/dashboard") && r.Method == http.MethodGet {
		GetSchoolDashboard(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ListSchools(w, r)
	case http.MethodPost:
		CreateSchool(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

