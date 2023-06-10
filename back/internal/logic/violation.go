package logic

import (
	"context"
	logicEntities "svalka-service/internal/logic/entities"
)

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
