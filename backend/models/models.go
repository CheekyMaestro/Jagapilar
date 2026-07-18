package models

import (
	"time"
)

// ==================== Database Models ====================

// User represents a parent or teacher
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	EmailContact string    `json:"email_contact"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`      // "parent", "teacher"
	SchoolID     *string   `json:"school_id,omitempty"` // nullable, for teachers
	CreatedAt    time.Time `json:"created_at"`
}

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

// Child represents a child's profile linked to a Neuro ID
type Child struct {
	ID        string    `json:"id"`
	NeuroID   string    `json:"neuro_id"`
	Name      string    `json:"name"`
	BirthYear int       `json:"birth_year"`
	Gender    string    `json:"gender"` // "L" or "P"
	CreatedBy string    `json:"created_by"` // user_id of the parent
	CreatedAt time.Time `json:"created_at"`
}

// TeacherChild mapping for teachers adding students to their dashboard
type TeacherChild struct {
	TeacherID string    `json:"teacher_id"`
	ChildID   string    `json:"child_id"`
	AddedAt   time.Time `json:"added_at"`
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
	UserID      *string    `json:"user_id,omitempty"` // null for students
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

type RegisterUserRequest struct {
	Role         string `json:"role"`
	Name         string `json:"name"`
	EmailContact string `json:"email_contact"`
	Password     string `json:"password"`
	SchoolName   string `json:"school_name,omitempty"` // we will just store this or ignore for now
}

type LoginRequest struct {
	EmailContact string `json:"email_contact"`
	Password     string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateChildRequest struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Gender    string `json:"gender"`
}

type AddChildToTeacherRequest struct {
	NeuroID string `json:"neuro_id"`
}

type SubmitResponseRequest struct {
	SessionID string `json:"session_id"`
	ItemCode  string `json:"item_code"`
	RawScore  int    `json:"raw_score"`
}

type SubmitAssessmentRequest struct {
	NeuroID    string           `json:"neuro_id"`
	Pillar     string           `json:"pillar"`
	TotalScore float64          `json:"total_score"`
	Zone       string           `json:"zone"`
	Items      []ItemSubmission `json:"items"`
}

type ItemSubmission struct {
	ItemCode string `json:"item_code"`
	RawScore int    `json:"raw_score"`
}

type AssessmentResult struct {
	ChildID       string        `json:"child_id"`
	PillarScores  []PillarScore `json:"pillar_scores"`
	Composite     *CompositeScore `json:"composite,omitempty"`
	Items         []ItemResponse `json:"items,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ChildDashboardItem struct {
	Child     Child           `json:"child"`
	Status    string          `json:"status"` // "pending", "parent_done", "complete"
	Composite *CompositeScore `json:"composite,omitempty"`
}
