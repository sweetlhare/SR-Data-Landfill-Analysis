package repository

import (
	"context"
	"errors"
	"fmt"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgconn"
)

// CreateAudit ...
func (rep repository) CreateAudit(ctx context.Context, auditRequest logicentities.AuditDbRequest) (*logicentities.AuditResposeAll, error) {

	// Audit ...
	builder := sq.Insert("audits").
		Columns("user_id", "survey_id", "date", "ai_generated_status").
		Values(
			auditRequest.UserID,
			auditRequest.SurveyID,
			auditRequest.Date,
			auditRequest.AiGeneratedStatus).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "audit.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	rows.Next()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	var auditID uint64
	err = rows.Scan(&auditID)
	if err != nil {
		return nil, err
	}

	if len(auditRequest.Violations) != 0 {
		// CreateAuditToViolations
		var auditsToViolations []repentities.AuditToViolation
		for _, v := range auditRequest.Violations {
			if v.Status == true {
				auditsToViolations = append(auditsToViolations, repentities.AuditToViolation{
					ViolationID: v.ID,
					AuditID:     auditID,
				})
			}
		}
		err = rep.CreateAuditToViolations(ctx, auditsToViolations)
		if err != nil {
			return nil, err
		}
	}

	return rep.GetAudit(ctx, auditID)
}

// GetAudit ...
func (rep repository) GetAudit(ctx context.Context, id uint64) (audits *logicentities.AuditResposeAll, err error) {
	// parse
	query := `
	select
			a.id,
			a.date,
			a.survey_id,
			a.ai_generated_status,
			COUNT(*) FILTER (where a."status" = TRUE) AS true_violations_count,
			json_build_object('id', u.id,'name', u."name",'position', u."position",'role', u."role",'phone', u.phone, 'email',u.email) as author,
			json_agg
			(
				json_build_object
				(
					'id', a.violation_id,
					'title', a.title,
					'default_status', a.default_status,
					'status', a.status
				)
			) as violations
			FROM
			(
				SELECT 
					a.*,
					v.id as violation_id, 
					v.description as "title",
					v.default_status,  
					v.id  = ANY(array_agg(atv.violation_id)) as status
				FROM audits_to_violations atv
				JOIN violations v on true
				inner join audits a on a.id = atv.audit_id 
				where a.id = $1
				GROUP BY a.id, v.id
				ORDER BY a.id, v.id
			) a
			LEFT join users u on u.id = a.user_id
			GROUP BY a.survey_id, a.id, a."date", a.ai_generated_status, u.id
			order by a."date"
	`

	q := pg.Query{
		Name:     "audits.GetById",
		QueryRaw: query,
	}

	var r repentities.Audit
	err = rep.client.PG().GetContext(ctx, &r, q, id)
	if err != nil {
		return nil, err
	}
	resp := rep.converter.Audit.ToLogic(r)
	return &resp, nil
}

// GetAudits ...
func (rep repository) GetAudits(ctx context.Context, surveyID uint64) (audits []logicentities.AuditResposeAll, err error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.Audit{})
	builder := squirrel.
		Select(columns...).
		From("audits").
		Where("survey_id = ?", surveyID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}

	var resultDB []repentities.Audit
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}

	for _, r := range resultDB { // допилить логику юзеров
		user, err := rep.GetUser(ctx, r.Author.ID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user = nil")
		}

		vs, count, err := rep.GetViolations(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		audits = append(audits, logicentities.AuditResposeAll{
			ID:                  r.ID,
			Date:                r.Date,
			AiGeneratedStatus:   r.AiGeneratedStatus,
			User:                *user,
			Violations:          vs,
			TrueViolationsCount: count,
		})
	}

	if len(audits) == 0 {
		audits = []logicentities.AuditResposeAll{}
	}

	return audits, nil
}

// CreateAuditToViolations ...
func (rep repository) CreateAuditToViolations(ctx context.Context, auditsToViolations []repentities.AuditToViolation) error {
	if len(auditsToViolations) == 0 {
		return nil
	}
	columns, _ := pg.ParseDbModel(repentities.AuditToViolation{})
	builder := sq.Insert("audits_to_violations").
		PlaceholderFormat(sq.Dollar).
		Columns(columns...)

	for _, auditToViolation := range auditsToViolations {
		_, volumes := pg.ParseDbModel(auditToViolation)
		builder = builder.Values(volumes...)
	}

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "audits_to_violations.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println(rows.Err())
	}
	err = rows.Err()
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			return errors.New(pgErr.Detail)
		}
		return err
	}

	return nil
}

// DeleteAudit ...
func (rep repository) DeleteAudit(ctx context.Context, id uint64) error {

	err := rep.DeleteViolationsByAuditID(ctx, id)
	if err != nil {
		return err
	}
	deleteBuilder := squirrel.Delete("audits").
		Where("id=?", id).PlaceholderFormat(sq.Dollar)

	// Generate the SQL and arguments
	query, v, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}
	q := pg.Query{
		Name:     "audits.Delete",
		QueryRaw: query,
	}
	// Execute the delete operation
	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}
