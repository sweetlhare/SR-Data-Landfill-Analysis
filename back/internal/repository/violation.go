package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

// CreateViolation ...
func (rep repository) CreateViolation(ctx context.Context, violation logicentities.Violation) error {
	violationDB := repentities.ViolationCreate{
		Description:   violation.Title,
		DefaultStatus: violation.DefaultStatus,
	}
	columns, values := pg.ParseDbModel(violationDB)
	builder := sq.Insert("violations").
		PlaceholderFormat(sq.Dollar).
		Columns(columns...).
		Values(values...)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "violation.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}

// GetViolations ...
func (rep repository) GetViolations(ctx context.Context, auditID uint64) (violations []logicentities.Violation, trueViolationCount uint64, err error) {
	// getAllviolations
	violations, err = rep.GetAllViolations(ctx)
	if err != nil {
		return nil, 0, err
	}

	// parse
	columns, _ := pg.ParseDbModel(repentities.AuditToViolation{})
	builder := squirrel.
		Select(columns...).
		From("audits_to_violations").
		Where("audit_id = ?", auditID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}

	var resultDB []*repentities.AuditToViolation
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, 0, err
	}
	trueViolations := make(map[uint64]interface{}, len(resultDB))
	for _, r := range resultDB {
		trueViolations[r.ViolationID] = 0
	}

	for i, v := range violations {
		_, ok := trueViolations[v.ID]
		if ok {
			v := violations[i]
			v.Status = true
			violations[i] = v
		}
	}

	return violations, uint64(len(trueViolations)), nil
}

// GetAllViolations ...
func (rep repository) GetAllViolations(ctx context.Context) ([]logicentities.Violation, error) {
	columns, _ := pg.ParseDbModel(repentities.ViolationResponseAll{})
	builder := sq.Select(columns...).
		From("violations").
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "violations.GetAll",
		QueryRaw: query,
	}

	var resultDB []*repentities.ViolationResponseAll
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}
	var result []logicentities.Violation
	for _, r := range resultDB {
		result = append(result, logicentities.Violation{
			ID:            r.ID,
			Title:         r.Title,
			DefaultStatus: r.DefaultStatus,
		})
	}
	return result, nil
}

// DeleteViolationsBySurveyID ...
func (rep repository) DeleteViolationsByAuditID(ctx context.Context, id uint64) error {
	deleteBuilder := squirrel.Delete("audits_to_violations").
		Where("audit_id=?", id).PlaceholderFormat(sq.Dollar)

	// Generate the SQL and arguments
	query, v, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}
	q := pg.Query{
		Name:     "violations.Delete",
		QueryRaw: query,
	}
	// Execute the delete operation
	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}
