package logicinterfaces

import logicentities "svalka-service/internal/logic/entities"

type AiClient interface {
	AnalyzeLandfill(cdnResult logicentities.CdnResult) (*logicentities.AiResult, error)
}
