package aiclient

import (
	"context"
	"svalka-service/internal/config"
	logicentities "svalka-service/internal/logic/entities"
	logicinterfaces "svalka-service/internal/logic/interfaces"
)

type aiClient struct {
	config config.AiClientConfig
}

func NewAiClient(_ context.Context) (logicinterfaces.AiClient, error) {
	return aiClient{
		config: AiClientConfig,
	}, nil
}

// AnalyzeLandfill ...
func (c aiClient) AnalyzeLandfill(cdnResult logicentities.CdnResult) (*logicentities.AiResult, error) {
	return &logicentities.AiResult{
		AiImagesPaths: []string{
			"test1_ai.png",
			"test2_ai.png",
		},
		Violations: []logicentities.Violation{},
	}, nil
}
