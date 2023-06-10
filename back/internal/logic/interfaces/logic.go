package logicinterfaces

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
)

// Logic ...
type Logic interface {
	// user
	CreateUser(ctx context.Context, user logicentities.UserCreate) error
	Login(ctx context.Context, login logicentities.Login) (*logicentities.Session, error)

	// Auth
	Auth(sessionID string, validRoles ...logicentities.UserRole) (*logicentities.SessionCredentials, error)

	// region
	CreateRegion(ctx context.Context, region logicentities.Region) error
	GetAllRegions(ctx context.Context) ([]logicentities.Region, error)

	// violation
	CreateViolation(ctx context.Context, violation logicentities.Violation) error
	GetAllViolations(ctx context.Context) ([]logicentities.Violation, error)

	// landfill
	CreateLandfill(ctx context.Context, landfill logicentities.LandfillCreate) (*logicentities.LandfillResponse, error)
	UpdateLandfill(ctx context.Context, landfill logicentities.LandfillUpdate) error
	GetLandfills(ctx context.Context, regionID uint64) ([]logicentities.LandfillAllResponse, error)
	GetLandfill(ctx context.Context, landfillID uint64) (*logicentities.LandfillResponse, error)
	DeleteLandfill(ctx context.Context, id uint64) error

	// survey
	CreateSurvey(ctx context.Context, surveyRequest logicentities.SurveyCreateRequest) (*logicentities.SurveyResponse, error)
	GetSurvey(ctx context.Context, surveyID uint64) (*logicentities.SurveyResponse, error)
	DeleteSurvey(ctx context.Context, id uint64) error

	// audit
	CreateAudit(ctx context.Context, auditRequest logicentities.AuditCreateRequest) (*logicentities.AuditResposeAll, error)
	DeleteAudit(ctx context.Context, id uint64) error
	GetAudit(ctx context.Context, id uint64) (*logicentities.AuditResposeAll, error)
}
