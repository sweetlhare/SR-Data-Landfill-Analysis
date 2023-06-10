package logicentities

type LandfillInfo struct {
	Illegal                 bool   `json:"illegal"`
	City                    string `json:"city"`
	Name                    string `json:"name"`
	Address                 string `json:"address"`
	Type                    string `json:"type"`
	Coordinates             string `json:"coordinates"  validate:"required"`
	CadastralNumber         string `json:"cadastral_number"`
	IllegalCadastralNumbers string `json:"illegal_cadastral_numbers"`
	CadastralCategory       string `json:"cadastral_category"`
	UsrArea                 string `json:"usr_area"`
	Area                    string `json:"area"`
	PreviewImagePath        string `json:"preview_image_path"`
}

type LandfillCreate struct {
	LandfillInfo
	RegionID uint64  `json:"region_id"  validate:"required"`
	Manager  Manager `json:"manager"`
}

type LandfillUpdate struct {
	ID uint64 `json:"id" validate:"required"`
	LandfillCreate
}

type LandfillAllResponse struct {
	ID uint64 `json:"id"`
	LandfillInfo
	Region          Region `json:"region"`
	ViolationsCount uint64 `json:"violations_count"`
}

type LandfillResponse struct {
	ID uint64 `json:"id"`
	LandfillInfo
	Region          Region              `json:"region"`
	Manager         Manager             `json:"manager"`
	ViolationsCount uint64              `json:"violations_count"`
	Violations      []Violation         `json:"violations"`
	Surveys         []SurveyResponseAll `json:"surveys"`
}
