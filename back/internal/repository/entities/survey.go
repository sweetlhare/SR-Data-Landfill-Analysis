package repentities

import (
	"database/sql"
	logicentities "svalka-service/internal/logic/entities"
	"time"
)

type SurveyResponseAll struct {
	ID              uint64    `json:"id"`
	LandfillID      uint64    `json:"landfill_id"`
	Date            time.Time `json:"date"`
	Author          User      `json:"author"`
	ViolationsCount uint64    `json:"violations_count"`
	Region          Region    `json:"region"`
}

type SurveyResponse struct {
	SurveyResponseAll
	RawImages []string `db:"raw_images"`
	AiImages  []string `db:"ai_images"`
	Audits    []Audit  `db:"audits"`
}

type SurveyCreate struct {
	UserID     uint64        `db:"user_id"`
	Date       time.Time     `db:"date"`
	LandfillID sql.NullInt64 `db:"landfill_id"`
}

// converter
type SurveyConverter struct {
	auditConverter  AuditConverter
	userConverter   UserConverter
	regionConverter RegionConverter
}

// ToLogic ...
func (c SurveyConverter) ToLogic(v SurveyResponse) logicentities.SurveyResponse {
	return logicentities.SurveyResponse{
		SurveyResponseAll: c.AllToLogic(v.SurveyResponseAll),
		RawImages:         v.RawImages,
		AiImages:          v.AiImages,
		Audits:            c.auditConverter.ToLogicArray(v.Audits...),
		Region:            c.regionConverter.ToLogic(&v.Region),
		ViolationsCount:   v.ViolationsCount,
	}
}

// AllToLogic ...
func (c SurveyConverter) AllToLogic(v SurveyResponseAll) logicentities.SurveyResponseAll {
	return logicentities.SurveyResponseAll{
		ID:              v.ID,
		LandfillID:      v.LandfillID,
		ViolationsCount: v.ViolationsCount,
		User:            c.userConverter.ToLogic(v.Author),
		Date:            v.Date,
	}
}

// AllToLogicArray ...
func (c SurveyConverter) AllToLogicArray(all ...SurveyResponseAll) (r []logicentities.SurveyResponseAll) {
	for _, a := range all {
		r = append(r, c.AllToLogic(a))
	}
	return r
}
