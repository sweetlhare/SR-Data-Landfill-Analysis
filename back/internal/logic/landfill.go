package logic

import (
	"context"
	"errors"
	logicEntities "svalka-service/internal/logic/entities"
)

// CreateLandfill ...
func (l logic) CreateLandfill(ctx context.Context, landfill logicEntities.LandfillCreate) (*logicEntities.LandfillResponse, error) {
	err := l.validator.CommonValidation(landfill)
	if err != nil {
		return nil, err
	}

	resp, err := l.rep.CreateLandfill(ctx, landfill)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteLandfill ...
func (l logic) DeleteLandfill(ctx context.Context, id uint64) error {
	return l.rep.DeleteLandfill(ctx, id)
}

// UpdateLandfill ...
func (l logic) UpdateLandfill(ctx context.Context, landfill logicEntities.LandfillUpdate) error {
	err := l.validator.CommonValidation(landfill)
	if err != nil {
		return err
	}

	err = l.rep.UpdateLandfill(ctx, landfill)
	if err != nil {
		return err
	}

	return nil
}

// GetLandfills ...
func (l logic) GetLandfills(ctx context.Context, regionID uint64) ([]logicEntities.LandfillAllResponse, error) {
	if regionID == 0 {
		return nil, errors.New("regionID = 0")
	}

	landfills, err := l.rep.GetAllLandfills(ctx, regionID)
	if err != nil {
		return nil, err
	}

	return landfills, nil
}

// GetLandfill ...
func (l logic) GetLandfill(ctx context.Context, landfillID uint64) (*logicEntities.LandfillResponse, error) {
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
