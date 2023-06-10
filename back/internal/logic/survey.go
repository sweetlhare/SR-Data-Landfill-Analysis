package logic

import (
	"context"
	"errors"
	logicEntities "svalka-service/internal/logic/entities"
	"time"
)

// CreateSurvey ...
func (l logic) CreateSurvey(ctx context.Context, s logicEntities.SurveyCreateRequest) (*logicEntities.SurveyResponse, error) {
	err := l.validator.CommonValidation(s)
	if err != nil {
		return nil, err
	}

	unique := l.rep.CheckUnique(ctx, s.Date, s.LandfillID)
	if !unique {
		return nil, errors.New("not unique")
	}

	// cdn
	cdnResult, err := l.cdnClient.SaveImages(int(s.MaxSize), s.RootImages...)
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

	// save survey
	surveyRequest := logicEntities.SurveyCreateDbRequest{
		UserID:     s.UserID,
		LandfillID: s.LandfillID,
		Date:       s.Date,
	}

	surveyID, err := l.rep.CreateSurvey(ctx, surveyRequest)
	if err != nil {
		return nil, err
	}
	if surveyID == 0 {
		return nil, errors.New("survey id = 0")
	}

	// save images
	var images []logicEntities.ImageCreate
	for _, i := range cdnResult.RawImagesPaths {
		images = append(images, logicEntities.ImageCreate{
			SurveyID:  surveyID,
			Path:      i,
			RawStatus: true,
		})
	}
	for _, i := range aiResult.AiImagesPaths {
		images = append(images, logicEntities.ImageCreate{
			SurveyID:  surveyID,
			Path:      i,
			RawStatus: false,
		})
	}
	err = l.rep.CreateImages(ctx, images...)
	if err != nil {
		return nil, err
	}

	// save audit
	_, err = l.rep.CreateAudit(ctx, logicEntities.AuditDbRequest{
		UserID:            s.UserID,
		Date:              time.Now(),
		AiGeneratedStatus: true,
		SurveyID:          surveyID,
		Violations:        aiResult.Violations,
	})
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

// DeleteSurvey ...
func (l logic) DeleteSurvey(ctx context.Context, id uint64) error {
	return l.rep.DeleteSurvey(ctx, id)
}

// GetSurvey ...
func (l logic) GetSurvey(ctx context.Context, surveyID uint64) (*logicEntities.SurveyResponse, error) {
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
