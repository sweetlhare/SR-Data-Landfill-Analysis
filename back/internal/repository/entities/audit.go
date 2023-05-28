package repentities

import (
	"time"
)

type Audit struct {
	ID                int64     `db:"id"`
	Date              time.Time `db:"date"`
	UserID            int64     `db:"user_id"`
	SurveyID          int64     `db:"survey_id"`
	AiGeneratedStatus bool      `db:"ai_generated_status"`
}

type AuditCreate struct {
	Date              time.Time `db:"date"`
	UserID            int64     `db:"user_id"`
	SurveyID          int64     `db:"survey_id"`
	AiGeneratedStatus bool      `db:"ai_generated_status"`
}
