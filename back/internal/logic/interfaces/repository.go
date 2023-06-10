package logicinterfaces

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	"time"
)

// Repository
type Repository interface {
	// user
	CreateUser(ctx context.Context, user logicentities.UserCreate) error
	GetUserByEmail(ctx context.Context, email string) (user *logicentities.UserWithPass, err error)

	// region
	CreateRegion(ctx context.Context, region logicentities.Region) error
	GetAllRegions(ctx context.Context) ([]logicentities.Region, error)

	// landfill
	CreateLandfill(ctx context.Context, landfill logicentities.LandfillCreate) (*logicentities.LandfillResponse, error)
	UpdateLandfill(ctx context.Context, landfill logicentities.LandfillUpdate) error
	GetAllLandfills(ctx context.Context, regionID uint64) ([]logicentities.LandfillAllResponse, error)
	GetLandfill(ctx context.Context, landfillID uint64) (*logicentities.LandfillResponse, error)
	DeleteLandfill(ctx context.Context, id uint64) error

	// violation
	CreateViolation(ctx context.Context, violation logicentities.Violation) error
	GetAllViolations(ctx context.Context) ([]logicentities.Violation, error)

	// survey
	CreateSurvey(ctx context.Context, surveyRequest logicentities.SurveyCreateDbRequest) (surveyID uint64, err error)
	GetSurvey(ctx context.Context, surveyID uint64) (*logicentities.SurveyResponse, error)
	DeleteSurvey(ctx context.Context, id uint64) error
	CheckUnique(ctx context.Context, date time.Time, landfillId uint64) bool

	// audit
	CreateAudit(ctx context.Context, auditRequest logicentities.AuditDbRequest) (*logicentities.AuditResposeAll, error)
	DeleteAudit(ctx context.Context, id uint64) error
	GetAudit(ctx context.Context, id uint64) (*logicentities.AuditResposeAll, error)

	// images
	CreateImages(ctx context.Context, images ...logicentities.ImageCreate) (err error)
}
