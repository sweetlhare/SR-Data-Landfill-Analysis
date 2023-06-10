package repentities

import (
	"database/sql"
	logicentities "svalka-service/internal/logic/entities"
)

type LandfillInfo struct {
	City                    sql.NullString `db:"city"`
	Name                    sql.NullString `db:"name"`
	Illegal                 bool           `db:"illegal"`
	Address                 sql.NullString `db:"address"`
	Coordinates             string         `db:"coordinates"`
	PreviewImagePath        sql.NullString `db:"preview_image_path"`
	Type                    sql.NullString `db:"type"`
	CadastralNumber         sql.NullString `db:"cadastral_number"`
	IllegalCadastralNumbers sql.NullString `db:"illegal_cadastral_numbers"`
	CadastralCategory       sql.NullString `db:"cadastral_category"`
	UsrArea                 sql.NullString `db:"usr_area"`
	Area                    sql.NullString `db:"area"`
}

type LandfillCreate struct {
	LandfillInfo
	RegionID uint64 `db:"region_id"`
	Manager  string `db:"manager"`
}

type LandfillResponseAll struct {
	ID uint64 `db:"id"`
	LandfillInfo
	Region          Region `db:"region"`
	ViolationsCount uint64 `db:"violations_count"`
}

type LandfillResponse struct {
	LandfillResponseAll
	Manager    Manager             `db:"manager"`
	Audit      Audit               `db:"audit"`
	Violations []Violation         `json:"violations"`
	Surveys    []SurveyResponseAll `db:"surveys"`
}

type LandfillConverter struct {
	surveyConverter    SurveyConverter
	managerConverter   ManagerConverter
	regionConverter    RegionConverter
	violationConverter ViolationConverter
}

// info to logic
func (c LandfillConverter) InfoToLogic(l LandfillInfo) logicentities.LandfillInfo {
	return logicentities.LandfillInfo{
		CadastralNumber:         l.CadastralNumber.String,
		IllegalCadastralNumbers: l.IllegalCadastralNumbers.String,
		CadastralCategory:       l.CadastralCategory.String,
		UsrArea:                 l.UsrArea.String,
		Area:                    l.Area.String,
		Illegal:                 l.Illegal,
		City:                    l.City.String,
		Name:                    l.Name.String,
		Address:                 l.Address.String,
		Type:                    l.Type.String,
		Coordinates:             l.Coordinates,
		PreviewImagePath:        l.PreviewImagePath.String,
	}
}

// info to logic
func (c LandfillConverter) InfoToDb(l logicentities.LandfillInfo) LandfillInfo {
	return LandfillInfo{
		CadastralNumber: sql.NullString{
			String: l.CadastralNumber,
			Valid:  true,
		},
		IllegalCadastralNumbers: sql.NullString{
			String: l.IllegalCadastralNumbers,
			Valid:  true,
		},
		CadastralCategory: sql.NullString{
			String: l.CadastralCategory,
			Valid:  true,
		},
		UsrArea: sql.NullString{
			String: l.UsrArea,
			Valid:  true,
		},
		Area: sql.NullString{
			String: l.Area,
			Valid:  true,
		},
		Illegal:     l.Illegal,
		Coordinates: l.Coordinates,
		City: sql.NullString{
			String: l.City,
			Valid:  true,
		},
		Name: sql.NullString{
			String: l.Name,
			Valid:  true,
		},
		Address: sql.NullString{
			String: l.Address,
			Valid:  true,
		},
		Type: sql.NullString{
			String: l.Type,
			Valid:  true,
		},
		PreviewImagePath: sql.NullString{
			String: l.PreviewImagePath,
			Valid:  true,
		},
	}
}

// info to logic
func (c LandfillConverter) AllToLogic(l LandfillResponseAll) logicentities.LandfillAllResponse {
	return logicentities.LandfillAllResponse{
		ID:              l.ID,
		LandfillInfo:    c.InfoToLogic(l.LandfillInfo),
		ViolationsCount: l.ViolationsCount,
		Region:          c.regionConverter.ToLogic(&l.Region),
	}
}

// AllToLogicArray
func (c LandfillConverter) AllToLogicArray(all ...LandfillResponseAll) (r []logicentities.LandfillAllResponse) {
	for _, a := range all {
		r = append(r, c.AllToLogic(a))
	}
	return r
}

// ToLogic ...
func (c LandfillConverter) ToLogic(l LandfillResponse) logicentities.LandfillResponse {
	r := logicentities.LandfillResponse{
		ID:              l.ID,
		LandfillInfo:    c.InfoToLogic(l.LandfillInfo),
		ViolationsCount: l.ViolationsCount,
		Region:          c.regionConverter.ToLogic(&l.Region),
		Manager:         c.managerConverter.ToLogic(l.Manager),
		Surveys:         c.surveyConverter.AllToLogicArray(l.Surveys...),
		Violations:      c.violationConverter.ToLogicArray(l.Violations...),
	}

	return r
}

// ToDB ...
func (c LandfillConverter) ToDB(l logicentities.LandfillCreate) LandfillCreate {
	r := LandfillCreate{
		LandfillInfo: c.InfoToDb(l.LandfillInfo),
		Manager:      c.managerConverter.ToDB(l.Manager),
		RegionID:     l.RegionID,
	}

	return r
}
