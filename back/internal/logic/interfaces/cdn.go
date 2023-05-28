package logicinterfaces

import logicentities "svalka-service/internal/logic/entities"

type CdnClient interface {
	SaveImages(images ...logicentities.File) (result *logicentities.CdnResult, err error)
}
