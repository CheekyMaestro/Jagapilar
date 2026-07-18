package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"jagapilar-backend/database"
	"jagapilar-backend/models"
)

// ==================== Scoring Engine ====================
// Implements the 360° scoring logic from the PRD:
// - Reverse-scoring for items O4, G1, M4 (formula: 6 - raw_score)
// - Sub-score per pillar: sum of 6 items → range 6-30
// - Composite score: average of 3 pillar sub-scores → range 6-30
// - Zone classification: Hijau (6.0-13.9), Kuning (14.0-21.9), Merah (22.0-30.0)
// - Convergence rule: if zones differ across pillars → needs_manual_review = true

// Reverse-scored item codes
var reverseItems = map[string]bool{
	"O4": true,
	"G1": true,
	"M4": true,
}

// CalculatePillarScore computes the sub-score for a single pillar
// applying reverse-scoring where needed. Range: 6-30.
func CalculatePillarScore(responses []models.ItemResponse) float64 {
	var total float64
	for _, r := range responses {
		if reverseItems[r.ItemCode] {
			total += float64(6 - r.RawScore)
		} else {
			total += float64(r.RawScore)
		}
	}
	return total
}

// ClassifyZone determines the risk zone based on a score (6-30)
func ClassifyZone(score float64) string {
	if score <= 13.9 {
		return "hijau"
	}
	if score <= 21.9 {
		return "kuning"
	}
	return "merah"
}

// CheckConvergence checks if all three pillar zones are the same
// If not, the assessment needs manual review (per De Los Reyes & Kazdin, 2005)
func CheckConvergence(zones []string) bool {
	if len(zones) < 2 {
		return false
	}
	first := zones[0]
	for _, z := range zones[1:] {
		if z != first {
			return true // needs manual review
		}
	}
	return false // all zones agree
}

// CalculateAndSavePillarScore calculates and persists a pillar score
func CalculateAndSavePillarScore(childID, pillar string, responses []models.ItemResponse) (*models.PillarScore, error) {
	subtotal := CalculatePillarScore(responses)
	zone := ClassifyZone(subtotal)

	query := `
		INSERT INTO pillar_scores (child_id, pillar, subtotal, zone, calculated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (child_id, pillar) DO UPDATE
		SET subtotal = $3, zone = $4, calculated_at = $5
	`
	now := time.Now()
	_, err := database.DB.Exec(query, childID, pillar, subtotal, zone, now)
	if err != nil {
		return nil, fmt.Errorf("failed to save pillar score: %w", err)
	}

	log.Printf("📊 Pillar score saved: child=%s pillar=%s subtotal=%.1f zone=%s", childID, pillar, subtotal, zone)

	return &models.PillarScore{
		ChildID:      childID,
		Pillar:       pillar,
		Subtotal:     subtotal,
		Zone:         zone,
		CalculatedAt: now,
	}, nil
}

// OnAllThreePillarsSubmitted is called when all 3 pillars have been submitted for a child.
// It calculates the composite 360° score and checks convergence.
func OnAllThreePillarsSubmitted(childID string) (*models.CompositeScore, error) {
	// Fetch all 3 pillar scores
	rows, err := database.DB.Query(
		"SELECT pillar, subtotal, zone FROM pillar_scores WHERE child_id = $1",
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pillar scores: %w", err)
	}
	defer rows.Close()

	var pillarScores []models.PillarScore
	var totalSubtotal float64
	var zones []string

	for rows.Next() {
		var ps models.PillarScore
		if err := rows.Scan(&ps.Pillar, &ps.Subtotal, &ps.Zone); err != nil {
			return nil, err
		}
		pillarScores = append(pillarScores, ps)
		totalSubtotal += ps.Subtotal
		zones = append(zones, ps.Zone)
	}

	if len(pillarScores) < 3 {
		return nil, fmt.Errorf("not all 3 pillars submitted yet (got %d)", len(pillarScores))
	}

	// Composite score = average of 3 pillars (equal weight 1:1:1)
	composite := totalSubtotal / 3.0
	compositeZone := ClassifyZone(composite)

	// Total raw = sum of all 3 pillar subtotals (range 18-90)
	totalRaw := totalSubtotal

	// Convergence check
	needsManualReview := CheckConvergence(zones)

	now := time.Now()
	query := `
		INSERT INTO composite_scores (child_id, composite, total_raw, composite_zone, needs_manual_review, calculated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (child_id) DO UPDATE
		SET composite = $2, total_raw = $3, composite_zone = $4, needs_manual_review = $5, calculated_at = $6
	`
	_, err = database.DB.Exec(query, childID, composite, totalRaw, compositeZone, needsManualReview, now)
	if err != nil {
		return nil, fmt.Errorf("failed to save composite score: %w", err)
	}

	result := &models.CompositeScore{
		ChildID:           childID,
		Composite:         composite,
		TotalRaw:          totalRaw,
		CompositeZone:     compositeZone,
		NeedsManualReview: needsManualReview,
		CalculatedAt:      now,
	}

	log.Printf("🎯 Composite 360° score: child=%s composite=%.1f zone=%s review=%v",
		childID, composite, compositeZone, needsManualReview)

	// Trigger referral if needed
	if compositeZone == "merah" || needsManualReview {
		if err := TriggerReferralWorkflow(childID, compositeZone, needsManualReview); err != nil {
			log.Printf("⚠️ Referral workflow failed: %v", err)
		}
	}

	return result, nil
}

// TriggerReferralWorkflow creates a referral record for high-risk children
func TriggerReferralWorkflow(childID, zone string, needsReview bool) error {
	referredTo := "BK_sekolah"
	if zone == "merah" {
		referredTo = "psikolog_klinis"
	}

	notes := fmt.Sprintf(`{"zone": "%s", "needs_review": %v, "auto_generated": true}`, zone, needsReview)

	query := `
		INSERT INTO referrals (child_id, referred_to, notes)
		VALUES ($1, $2, $3::jsonb)
	`
	_, err := database.DB.Exec(query, childID, referredTo, notes)
	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}

	log.Printf("🚨 Referral created: child=%s → %s", childID, referredTo)
	return nil
}

// GetAssessmentResult retrieves the full assessment result for a child
func GetAssessmentResult(childID string) (*models.AssessmentResult, error) {
	result := &models.AssessmentResult{
		ChildID: childID,
	}

	// Fetch pillar scores
	rows, err := database.DB.Query(
		"SELECT child_id, pillar, subtotal, zone, calculated_at FROM pillar_scores WHERE child_id = $1",
		childID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ps models.PillarScore
		if err := rows.Scan(&ps.ChildID, &ps.Pillar, &ps.Subtotal, &ps.Zone, &ps.CalculatedAt); err != nil {
			return nil, err
		}
		result.PillarScores = append(result.PillarScores, ps)
	}

	// Fetch composite score
	var cs models.CompositeScore
	err = database.DB.QueryRow(
		"SELECT child_id, composite, total_raw, composite_zone, needs_manual_review, calculated_at FROM composite_scores WHERE child_id = $1",
		childID,
	).Scan(&cs.ChildID, &cs.Composite, &cs.TotalRaw, &cs.CompositeZone, &cs.NeedsManualReview, &cs.CalculatedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == nil {
		result.Composite = &cs
	}

	return result, nil
}

// SeedAssessmentItems inserts the 18 assessment items if they don't exist
func SeedAssessmentItems() error {
	items := []models.AssessmentItem{
		// Pilar Orang Tua (O1-O6)
		{Code: "O1", Pillar: "parent", QuestionText: "Anak saya tantrum/marah berlebihan saat gadget diambil atau dibatasi waktunya.", IsReverse: false, Construct: "Dopamine withdrawal", ReferenceSource: "Firth et al. (2019)"},
		{Code: "O2", Pillar: "parent", QuestionText: "Anak saya terus meminta tambahan waktu bermain gadget meski sudah diberi batas waktu yang jelas.", IsReverse: false, Construct: "Toleransi & eskalasi", ReferenceSource: "Firth et al. (2019)"},
		{Code: "O3", Pillar: "parent", QuestionText: "Anak saya sulit tidur atau tidur larut malam karena bermain gadget.", IsReverse: false, Construct: "Gangguan tidur", ReferenceSource: "Hutton et al. (2020)"},
		{Code: "O4", Pillar: "parent", QuestionText: "Anak saya masih bisa beraktivitas tanpa gadget selama lebih dari 60 menit tanpa rewel.", IsReverse: true, Construct: "Delayed gratification", ReferenceSource: "Deci & Ryan (2000)"},
		{Code: "O5", Pillar: "parent", QuestionText: "Anak saya lebih memilih main gadget sendirian daripada bermain/mengobrol dengan keluarga atau teman.", IsReverse: false, Construct: "Experiential avoidance", ReferenceSource: "Hayes et al. (1996)"},
		{Code: "O6", Pillar: "parent", QuestionText: "Saya sendiri sering bermain gadget saat sedang menemani atau mengawasi anak.", IsReverse: false, Construct: "Parental digital neglect", ReferenceSource: "Radesky & Christakis (2016)"},

		// Pilar Guru (G1-G6)
		{Code: "G1", Pillar: "teacher", QuestionText: "Anak ini mampu fokus mengerjakan tugas/mendengarkan penjelasan >15 menit tanpa teralihkan.", IsReverse: true, Construct: "Sustained attention", ReferenceSource: "Shou et al. (2025)"},
		{Code: "G2", Pillar: "teacher", QuestionText: "Anak ini tampak lelah, mengantuk, atau kurang bersemangat di kelas pada pagi hari.", IsReverse: false, Construct: "Indikasi gangguan tidur", ReferenceSource: "Hutton et al. (2020)"},
		{Code: "G3", Pillar: "teacher", QuestionText: "Anak ini menunjukkan impulsivitas — sulit menunggu giliran, memotong pembicaraan, bertindak tanpa berpikir.", IsReverse: false, Construct: "Impulsivitas / ADHD-like", ReferenceSource: "Shou et al. (2025)"},
		{Code: "G4", Pillar: "teacher", QuestionText: "Anak ini kesulitan mengingat instruksi atau materi yang baru saja disampaikan.", IsReverse: false, Construct: "Working memory", ReferenceSource: "Nagata et al. (2024)"},
		{Code: "G5", Pillar: "teacher", QuestionText: "Anak ini menarik diri dari interaksi sosial dengan teman sebaya saat jam istirahat.", IsReverse: false, Construct: "Isolasi sosial", ReferenceSource: "Deci & Ryan (2000)"},
		{Code: "G6", Pillar: "teacher", QuestionText: "Prestasi/performa akademik anak ini menurun dibanding awal semester tanpa sebab akademik yang jelas.", IsReverse: false, Construct: "Penurunan performa akademik", ReferenceSource: "National Academies (2024)"},

		// Pilar Murid (M1-M6)
		{Code: "M1", Pillar: "student", QuestionText: "Aku jadi gelisah atau nggak enak rasanya, kalau lama nggak pegang HP/gadget.", IsReverse: false, Construct: "Craving / withdrawal subjektif", ReferenceSource: "Firth et al. (2019)"},
		{Code: "M2", Pillar: "student", QuestionText: "Awalnya aku cuma mau main sebentar, tapi ujung-ujungnya main gadget lama banget.", IsReverse: false, Construct: "Loss of control", ReferenceSource: "Firth et al. (2019)"},
		{Code: "M3", Pillar: "student", QuestionText: "Aku main gadget biar bisa lupa, kalau lagi sedih, bosan, atau marah.", IsReverse: false, Construct: "Experiential avoidance", ReferenceSource: "Hayes et al. (1996)"},
		{Code: "M4", Pillar: "student", QuestionText: "Aku masih suka main di luar rumah atau ngobrol langsung sama teman dan keluarga, dibanding main gadget terus-terusan.", IsReverse: true, Construct: "Restorative engagement", ReferenceSource: "Kaplan & Kaplan (1989)"},
		{Code: "M5", Pillar: "student", QuestionText: "Aku suka kepikiran gadget terus, padahal lagi belajar atau ngerjain PR.", IsReverse: false, Construct: "Intrusive thoughts", ReferenceSource: "Firth et al. (2019)"},
		{Code: "M6", Pillar: "student", QuestionText: "Kalau lagi ada masalah, aku lebih milih main gadget sendirian daripada cerita ke orang tua atau guru.", IsReverse: false, Construct: "Self-isolation", ReferenceSource: "Deci & Ryan (2000)"},
	}

	for _, item := range items {
		query := `
			INSERT INTO assessment_items (code, pillar, question_text, is_reverse, construct, reference_source)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (code) DO NOTHING
		`
		_, err := database.DB.Exec(query, item.Code, item.Pillar, item.QuestionText, item.IsReverse, item.Construct, item.ReferenceSource)
		if err != nil {
			return fmt.Errorf("failed to seed item %s: %w", item.Code, err)
		}
	}

	log.Println("✅ 18 assessment items seeded successfully")
	return nil
}
