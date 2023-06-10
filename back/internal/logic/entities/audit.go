package logicentities

import "time"

type AuditResposeAll struct {
	ID                  uint64      `json:"id"`
	Date                time.Time   `json:"date"`
	AiGeneratedStatus   bool        `json:"ai_generated_status"`
	User                User        `json:"author"`
	TrueViolationsCount uint64      `json:"true_violations_count"`
	Violations          []Violation `json:"violations"`
}

type AuditCreateRequest struct {
	UserID     uint64      `json:"author_id"`
	SurveyID   uint64      `json:"survey_id"  validate:"required"`
	Violations []Violation `json:"violations"`
}

type AuditDbRequest struct {
	UserID            uint64
	Date              time.Time
	AiGeneratedStatus bool
	SurveyID          uint64
	Violations        []Violation
}
