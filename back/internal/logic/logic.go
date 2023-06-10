package logic

import (
	"context"
	"errors"
	logicInterfaces "svalka-service/internal/logic/interfaces"
	validator "svalka-service/internal/logic/validator"
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
