package logic

import (
	"context"
	"errors"
	logicEntities "svalka-service/internal/logic/entities"
	logicInterfaces "svalka-service/internal/logic/interfaces"
	validator "svalka-service/internal/logic/validator"
	"time"
)

// logic ...
type logic struct {
	rep       logicInterfaces.Repository
	sessRep   logicInterfaces.SessionRepository
	aiClient  logicInterfaces.AiClient
	cdnClient logicInterfaces.CdnClient
	validator logicInterfaces.Validator
}

// NewLogic...
func NewLogic(
	ctx context.Context,
	rep logicInterfaces.Repository,
	sessRep logicInterfaces.SessionRepository,
	aiClient logicInterfaces.AiClient,
	cdnClient logicInterfaces.CdnClient,
) (logicInterfaces.Logic, error) {
	if rep == nil {
		return nil, errors.New("repository not initialized")
	}
	if sessRep == nil {
		return nil, errors.New("session repository not initialized")
	}
	if aiClient == nil {
		return nil, errors.New("ai clinent not initialized")
	}
	if cdnClient == nil {
		return nil, errors.New("cdn clinent not initialized")
	}
	return logic{
		rep:       rep,
		sessRep:   sessRep,
		aiClient:  aiClient,
		cdnClient: cdnClient,
		validator: validator.NewValidator(),
	}, nil
}

// CreateUser ...
func (l logic) CreateUser(ctx context.Context, user logicEntities.User) error {
	err := l.validator.CommonValidation(user)
	if err != nil {
		return err
	}
	err = l.rep.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// Login ...
func (l logic) Login(ctx context.Context, login logicEntities.Login) (*logicEntities.Session, error) {
	err := l.validator.CommonValidation(login)
	if err != nil {
		return nil, err
	}

	ok, err := l.rep.Login(ctx, login)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("user does not exist")
	}

	session, err := l.sessRep.CreateSession(login.Email)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// CreateRegion ...
func (l logic) CreateRegion(ctx context.Context, region logicEntities.Region) error {
	err := l.validator.CommonValidation(region)
	if err != nil {
		return err
	}

	err = l.rep.CreateRegion(ctx, region)
	if err != nil {
		return err
	}

	return nil
}

// GetAllRegions ...
func (l logic) GetAllRegions(ctx context.Context) ([]logicEntities.Region, error) {
	regions, err := l.rep.GetAllRegions(ctx)
	if err != nil {
		return nil, err
	}

	return regions, nil
}

// CreateLandfill ...
func (l logic) CreateLandfill(ctx context.Context, landfill logicEntities.Landfill) error {
	err := l.validator.CommonValidation(landfill)
	if err != nil {
		return err
	}

	err = l.rep.CreateLandfill(ctx, landfill)
	if err != nil {
		return err
	}

	return nil
}

// GetLandfills ...
func (l logic) GetLandfills(ctx context.Context, regionID int64) ([]logicEntities.LandfillAllResponse, error) {
	if regionID == 0 {
		return nil, errors.New("regionID = 0")
	}

	landfills, err := l.rep.GetLandfills(ctx, regionID)
	if err != nil {
		return nil, err
	}

	return landfills, nil
}

// CreateViolation ...
func (l logic) CreateViolation(ctx context.Context, violation logicEntities.Violation) error {
	err := l.validator.CommonValidation(violation)
	if err != nil {
		return err
	}

	err = l.rep.CreateViolation(ctx, violation)
	if err != nil {
		return err
	}

	return nil
}

// GetAllViolations ...
func (l logic) GetAllViolations(ctx context.Context) ([]logicEntities.Violation, error) {
	violations, err := l.rep.GetAllViolations(ctx)
	if err != nil {
		return nil, err
	}

	return violations, nil
}

// CreateSurvey ...
func (l logic) CreateSurvey(ctx context.Context, surveyRootRequest logicEntities.SurveyRootRequest) (*logicEntities.Survey, error) {
	err := l.validator.CommonValidation(surveyRootRequest)
	if err != nil {
		return nil, err
	}
	// cdn
	cdnResult, err := l.cdnClient.SaveImages(surveyRootRequest.RootImages...)
	if err != nil {
		return nil, err
	}
	if cdnResult == nil {
		return nil, errors.New("cdn result = nil")
	}

	// ai
	aiResult, err := l.aiClient.AnalyzeLandfill(*cdnResult)
	if err != nil {
		return nil, err
	}
	if aiResult == nil {
		return nil, errors.New("cdn ai result = nil")
	}

	// create survey
	surveyRequest := logicEntities.SurveyRequest{
		LandfillID: surveyRootRequest.LandfillID,
		Date:       surveyRootRequest.Date,
	}

	surveyID, err := l.rep.CreateSurvey(ctx, surveyRequest)
	if err != nil {
		return nil, err
	}
	if surveyID == 0 {
		return nil, errors.New("survey id = 0")
	}

	// save processing result
	surveyResult := logicEntities.SurveyResult{
		SurveyID:  surveyID,
		CdnResult: *cdnResult,
		AiResult:  *aiResult,
	}
	err = l.rep.SaveSurveyResult(ctx, surveyResult)
	if err != nil {
		return nil, err
	}

	// get survey
	survey, err := l.GetSurvey(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	return survey, err
}

// GetSurvey ...
func (l logic) GetSurvey(ctx context.Context, surveyID int64) (*logicEntities.Survey, error) {
	if surveyID == 0 {
		return nil, errors.New("surveyID = 0")
	}
	survey, err := l.rep.GetSurvey(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("failed get survey from rep")
	}
	return survey, nil
}

//  CreateAudit ...
func (l logic) CreateAudit(ctx context.Context, auditRequest logicEntities.AuditRequest) (*logicEntities.Survey, error) {
	auditRequest.Date = time.Now()
	err := l.validator.CommonValidation(auditRequest)
	if err != nil {
		return nil, err
	}

	// get email by session
	email, err := l.sessRep.GetValue(auditRequest.SessionID)
	if err != nil {
		return nil, err
	}

	// create audit
	auditRequest.Email = email
	auditRequest.Date = time.Now()
	err = l.rep.CreateAudit(ctx, auditRequest)
	if err != nil {
		return nil, err
	}

	// get survey
	survey, err := l.GetSurvey(ctx, auditRequest.SurveyID)
	if err != nil {
		return nil, err
	}
	return survey, err
}

// GetLandfill ...
func (l logic) GetLandfill(ctx context.Context, landfillID int64) (*logicEntities.LandfillExtended, error) {
	if landfillID == 0 {
		return nil, errors.New("landfillID = 0")
	}

	// get landfill
	landfill, err := l.rep.GetLandfill(ctx, landfillID)
	if err != nil {
		return nil, err
	}

	return landfill, err
}
