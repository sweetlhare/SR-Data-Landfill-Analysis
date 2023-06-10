package repentities

import (
	"database/sql"
	logicentities "svalka-service/internal/logic/entities"
	"time"
)

type Audit struct {
	ID                  uint64      `json:"id"`
	Date                time.Time   `json:"date"`
	SurveyID            uint64      `json:"survey_id"`
	AiGeneratedStatus   bool        `json:"ai_generated_status"`
	Author              User        `json:"author"`
	TrueViolationsCount uint64      `json:"true_violations_count"`
	Violations          []Violation `json:"violations"`
}

type AuditCreate struct {
	Date              time.Time     `db:"date"`
	UserID            sql.NullInt64 `db:"user_id"`
	SurveyID          uint64        `db:"survey_id"`
	AiGeneratedStatus bool          `db:"ai_generated_status"`
}

type AuditConverter struct {
	userConverter      UserConverter
	violationConverter ViolationConverter
}

// ToLogic ...
func (c AuditConverter) ToLogic(a Audit) logicentities.AuditResposeAll {
	return logicentities.AuditResposeAll{
		ID:                  a.ID,
		Date:                a.Date,
		AiGeneratedStatus:   a.AiGeneratedStatus,
		User:                c.userConverter.ToLogic(a.Author),
		Violations:          c.violationConverter.ToLogicArray(a.Violations...),
		TrueViolationsCount: a.TrueViolationsCount,
	}
}

// ToLogic ...
func (c AuditConverter) ToLogicArray(vs ...Audit) (r []logicentities.AuditResposeAll) {
	for _, v := range vs {
		r = append(r, c.ToLogic(v))
	}
	return r
}

// ToLogic ...
func (c AuditConverter) ToDB(auditRequest logicentities.AuditDbRequest) AuditCreate {
	var (
		userID            sql.NullInt64
		aiGeneratedStatus bool
	)
	if auditRequest.UserID == 0 { // перенести в
		userID = sql.NullInt64{
			Valid: true,
			Int64: int64(auditRequest.UserID),
		}
		aiGeneratedStatus = true
	}
	return AuditCreate{
		Date:              time.Now(),
		UserID:            userID,
		SurveyID:          auditRequest.SurveyID,
		AiGeneratedStatus: aiGeneratedStatus,
	}
}
