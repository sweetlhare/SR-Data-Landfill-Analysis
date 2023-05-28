package logicinterfaces

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
)

// Logic ...
type Logic interface {
	CreateUser(ctx context.Context, user logicentities.User) error
	Login(ctx context.Context, login logicentities.Login) (*logicentities.Session, error)
	CreateRegion(ctx context.Context, region logicentities.Region) error
	GetAllRegions(ctx context.Context) ([]logicentities.Region, error)
	CreateLandfill(ctx context.Context, landfill logicentities.Landfill) error
	GetLandfills(ctx context.Context, regionID int64) ([]logicentities.LandfillAllResponse, error)
	CreateViolation(ctx context.Context, violation logicentities.Violation) error
	GetAllViolations(ctx context.Context) ([]logicentities.Violation, error)
	CreateSurvey(ctx context.Context, surveyRequest logicentities.SurveyRootRequest) (*logicentities.Survey, error)
	GetSurvey(ctx context.Context, surveyID int64) (*logicentities.Survey, error)
	CreateAudit(ctx context.Context, auditRequest logicentities.AuditRequest) (*logicentities.Survey, error)
	GetLandfill(ctx context.Context, landfillID int64) (*logicentities.LandfillExtended, error)
}
