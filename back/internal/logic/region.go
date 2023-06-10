package logic

import (
	"context"
	logicEntities "svalka-service/internal/logic/entities"
)

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
