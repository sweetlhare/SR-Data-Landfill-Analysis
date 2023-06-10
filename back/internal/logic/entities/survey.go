package logicentities

import "time"

type SurveyCreateRequest struct {
	UserID     uint64    `json:"user_id"`
	LandfillID uint64    `json:"landfill_id"`
	MaxSize    uint64    `json:"max_size"`
	Date       time.Time `json:"date"`
	RootImages []File    `json:""`
}

type SurveyCreateDbRequest struct {
	UserID     uint64    `json:"user_id"`
	LandfillID uint64    `json:"landfill_id"`
	Date       time.Time `json:"date"`
}

type SurveyResponseAll struct {
	ID              uint64    `json:"id"`
	LandfillID      uint64    `json:"landfill_id"`
	ViolationsCount uint64    `json:"violations_count"`
	User            User      `json:"author"`
	Date            time.Time `json:"date" validate:"required"`
}

type SurveyResponse struct {
	SurveyResponseAll
	Region          Region            `json:"region"`
	RawImages       []string          `json:"raw_images"`
	AiImages        []string          `json:"ai_images"`
	ViolationsCount uint64            `json:"violations_count"`
	Audits          []AuditResposeAll `json:"audits"`
}
