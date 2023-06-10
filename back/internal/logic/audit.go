package logic

import (
	"context"
	logicEntities "svalka-service/internal/logic/entities"
	"time"
)

// CreateAudit ...
func (l logic) CreateAudit(ctx context.Context, auditRequest logicEntities.AuditCreateRequest) (*logicEntities.AuditResposeAll, error) {
	err := l.validator.CommonValidation(auditRequest)
	if err != nil {
		return nil, err
	}

	auditDbRequest := logicEntities.AuditDbRequest{
		UserID:            auditRequest.UserID,
		AiGeneratedStatus: false, // by user
		SurveyID:          auditRequest.SurveyID,
		Date:              time.Now(),
		Violations:        auditRequest.Violations,
	}

	// create db audit
	resp, err := l.rep.CreateAudit(ctx, auditDbRequest)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteAudit ...
func (l logic) DeleteAudit(ctx context.Context, id uint64) error {
	return l.rep.DeleteAudit(ctx, id)
}

// GetAudit ...
func (l logic) GetAudit(ctx context.Context, id uint64) (*logicEntities.AuditResposeAll, error) {
	return l.rep.GetAudit(ctx, id)
}
