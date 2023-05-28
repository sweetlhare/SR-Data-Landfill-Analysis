package repentities

import (
	"time"
)

type Survey struct {
	ID         int64     `db:"id"`
	Date       time.Time `db:"date"`
	LandfillID int64     `db:"landfill_id"`
}

type SurveyCreate struct {
	Date       time.Time `db:"date"`
	LandfillID int64     `db:"landfill_id"`
}
