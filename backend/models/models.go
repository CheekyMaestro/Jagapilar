package models

import (
	"time"
)

// ==================== Database Models ====================

// School represents a registered school
type School struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	City                  string    `json:"city"`
	GradeLevel            string    `json:"grade_level"`
	PrincipalName         string    `json:"principal_name"`
	TotalClasses          int       `json:"total_classes"`
	PrincipalConsentDocURL string   `json:"principal_consent_doc_url,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

// Child represents an anonymized child record
type Child struct {
	ID        string    `json:"id"`
	AnonCode  string    `json:"anon_code"`
	SchoolID  string    `json:"school_id"`
	Grade     int       `json:"grade"`
	BirthYear int       `json:"birth_year"`
	CreatedAt time.Time `json:"created_at"`
}

// Informant represents a parent/teacher/student linked to a child
type Informant struct {
	ID              string     `json:"id"`
	ChildID         string     `json:"child_id"`
	Role            string     `json:"role"` // "parent", "teacher", "student"
	ContactHash     string     `json:"contact_hash,omitempty"`
	AccessToken     string     `json:"access_token"`
	TokenExpiresAt  *time.Time `json:"token_expires_at,omitempty"`
	ConsentSignedAt *time.Time `json:"consent_signed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// AssessmentItem represents one of the 18 questionnaire items
type AssessmentItem struct {
	Code            string `json:"code"`
	Pillar          string `json:"pillar"`
	QuestionText    string `json:"question_text"`
	IsReverse       bool   `json:"is_reverse"`
	Construct       string `json:"construct"`
	ReferenceSource string `json:"reference_source,omitempty"`
}

// AssessmentSession represents a single assessment session
type AssessmentSession struct {
	ID          string     `json:"id"`
	ChildID     string     `json:"child_id"`
	InformantID string     `json:"informant_id"`
	Pillar      string     `json:"pillar"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Status      string     `json:"status"` // "draft", "submitted"
	CreatedAt   time.Time  `json:"created_at"`
}

// ItemResponse represents a single answer to one assessment item
type ItemResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	ItemCode  string    `json:"item_code"`
	RawScore  int       `json:"raw_score"` // 1-5
	SavedAt   time.Time `json:"saved_at"`
}

// PillarScore represents the calculated sub-score for one pillar
type PillarScore struct {
	ChildID      string    `json:"child_id"`
	Pillar       string    `json:"pillar"`
	Subtotal     float64   `json:"subtotal"` // 6-30
	Zone         string    `json:"zone"`     // "hijau", "kuning", "merah"
	CalculatedAt time.Time `json:"calculated_at"`
}

// CompositeScore represents the 360° composite result
type CompositeScore struct {
	ChildID           string    `json:"child_id"`
	Composite         float64   `json:"composite"`      // average of 3 pillars, 6-30
	TotalRaw          float64   `json:"total_raw"`       // sum of 18 items, 18-90
	CompositeZone     string    `json:"composite_zone"`
	NeedsManualReview bool      `json:"needs_manual_review"`
	CalculatedAt      time.Time `json:"calculated_at"`
}

// ==================== Request/Response DTOs ====================

// CreateSchoolRequest is the request body for school registration
type CreateSchoolRequest struct {
	Name          string `json:"name"`
	City          string `json:"city"`
	GradeLevel    string `json:"grade_level"`
	PrincipalName string `json:"principal_name"`
	TotalClasses  int    `json:"total_classes"`
}

// CreateChildRequest is the request body for registering a child
type CreateChildRequest struct {
	SchoolID  string `json:"school_id"`
	Grade     int    `json:"grade"`
	BirthYear int    `json:"birth_year"`
}

// CreateInformantRequest is the request body for creating an informant
type CreateInformantRequest struct {
	Role        string `json:"role"`
	ContactHash string `json:"contact_hash,omitempty"`
}

// SubmitResponseRequest is the request body for submitting a single item response
type SubmitResponseRequest struct {
	SessionID string `json:"session_id"`
	ItemCode  string `json:"item_code"`
	RawScore  int    `json:"raw_score"`
}

// SubmitAssessmentRequest is the request body for submitting a full assessment
type SubmitAssessmentRequest struct {
	Pillar     string               `json:"pillar"`
	TotalScore float64              `json:"total_score"`
	Zone       string               `json:"zone"`
	Items      []ItemSubmission     `json:"items"`
	ChildID    string               `json:"child_id,omitempty"`
}

// ItemSubmission represents a single item answer in the submission
type ItemSubmission struct {
	ItemCode string `json:"item_code"`
	RawScore int    `json:"raw_score"`
}

// AssessmentResult is the response for assessment results
type AssessmentResult struct {
	ChildID       string        `json:"child_id"`
	PillarScores  []PillarScore `json:"pillar_scores"`
	Composite     *CompositeScore `json:"composite,omitempty"`
	Items         []ItemResponse `json:"items,omitempty"`
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TokenValidationRequest is the request to validate an access token
type TokenValidationRequest struct {
	Token string `json:"token"`
}

// TokenValidationResponse is the response from token validation
type TokenValidationResponse struct {
	Valid     bool   `json:"valid"`
	Role      string `json:"role,omitempty"`
	ChildID   string `json:"child_id,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// RegisterChildFullRequest is the request body for registering a child with 3 informants
type RegisterChildFullRequest struct {
	SchoolID  string `json:"school_id"`
	Grade     int    `json:"grade"`
	BirthYear int    `json:"birth_year"`
}

// ChildWithTokens represents a child and their 3 magic links
type ChildWithTokens struct {
	Child        Child  `json:"child"`
	ParentToken  string `json:"parent_token"`
	TeacherToken string `json:"teacher_token"`
	StudentToken string `json:"student_token"`
}

// SchoolDashboardResponse represents the dashboard data
type SchoolDashboardResponse struct {
	School      School               `json:"school"`
	TotalHijau  int                  `json:"total_hijau"`
	TotalKuning int                  `json:"total_kuning"`
	TotalMerah  int                  `json:"total_merah"`
	TotalReview int                  `json:"total_review"`
	Children    []ChildDashboardItem `json:"children"`
}

type ChildDashboardItem struct {
	Child        Child           `json:"child"`
	Composite    *CompositeScore `json:"composite,omitempty"`
	ParentToken  string          `json:"parent_token,omitempty"`
	TeacherToken string          `json:"teacher_token,omitempty"`
	StudentToken string          `json:"student_token,omitempty"`
}
