package logicentities

import (
	"mime/multipart"
	"time"
)

type UserInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email" validate:"required"`
}

type User struct {
	UserInfo
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type Region struct {
	ID    int64  `json:"id"`
	Title string `json:"title"  validate:"required"`
}

type Manager struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type LandfillInfo struct {
	ID               int64  `json:"id"`
	RegionID         int64  `json:"region_id"  validate:"required"`
	Illegal          bool   `json:"illegal"`
	City             string `json:"city"`
	Address          string `json:"address"`
	Type             string `json:"type"`
	Coordinates      string `json:"coordinates"  validate:"required"`
	PreviewImagePath string `json:"preview_image_path"`
}

type Landfill struct {
	LandfillInfo
	Manager Manager `json:"manager"`
}

type LandfillAllResponse struct {
	LandfillInfo
	ViolationsCount int64 `json:"violations_count"`
}

type Violation struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Default bool   `json:"default"`
	Status  bool   `json:"status"`
}

type Audit struct {
	ID                  int64       `json:"id"`
	Date                time.Time   `json:"date"`
	AiGeneratedStatus   bool        `json:"ai_generated_status"`
	Auditor             UserInfo    `json:"auditor"`
	Violations          []Violation `json:"violations"`
	TrueViolationsCount int         `json:"true_violations_count"`
}

type AuditRequest struct {
	SessionID         string      `json:""` // TODO: перенести в контекст
	Date              time.Time   `json:"date"  validate:"required"`
	AiGeneratedStatus bool        `json:"ai_generated_status"`
	Email             string      `json:"email"`
	SurveyID          int64       `json:"survey_id"  validate:"required"`
	Violations        []Violation `json:"violations"`
}

type File struct {
	*multipart.FileHeader
}

type SurveyRootRequest struct {
	LandfillID int64     `json:"landfill_id"`
	Date       time.Time `json:"date"`
	RootImages []File    `json:""`
}

type CdnResult struct {
	RawImagesPaths []string
}

type AiResult struct {
	AiImagesPaths []string
	Violations    []Violation
}

type SurveyRequest struct {
	LandfillID int64
	Date       time.Time
}

type SurveyResult struct {
	SurveyID  int64
	CdnResult CdnResult
	AiResult  AiResult
}

type Survey struct {
	ID        int64     `json:"id"`
	Date      time.Time `json:"date" validate:"required"`
	RawImages []string  `json:"raw_images"`
	AiImages  []string  `json:"ai_images"`
	Audits    []Audit   `json:"audits"`
}

type LandfillExtended struct {
	LandfillInfo
	ViolationsCount int64    `json:"violations_count"`
	Manager         Manager  `json:"manager"`
	Surveys         []Survey `json:"surveys"`
}

type Session struct {
	ID string `json:"session_id"`
}
