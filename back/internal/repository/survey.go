package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

// CreateSurvey ...
func (rep repository) CreateSurvey(ctx context.Context, surveyRequest logicentities.SurveyCreateDbRequest) (surveyID uint64, err error) {
	builder := sq.Insert("surveys").
		PlaceholderFormat(sq.Dollar).
		Columns("date", "landfill_id", "user_id").
		Values(surveyRequest.Date, surveyRequest.LandfillID, surveyRequest.UserID).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := pg.Query{
		Name:     "survey.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rows.Next()
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

// GetSurvey ...
func (rep repository) GetSurvey(ctx context.Context, surveyID uint64) (*logicentities.SurveyResponse, error) {
	query := `
	select
	s.id, 
	s."date", 
	s.landfill_id,
	i.raw_images,
	i.ai_images,
	jsonb_build_object('id', u2.id, 'name', u2."name", 'position', u2."position",'role', u2."role", 'phone', u2.phone, 'email', u2.email) as author,
	jsonb_build_object('id', r.id, 'name', r."name") as region,
	json_agg(a.*) as audits
	FROM public.surveys s
	LEFT JOIN 
		(
			select
			a.id,
			TO_CHAR(a."date", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as date,
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
				GROUP BY a.id, v.id
				ORDER BY a.id, v.id
			) a
			LEFT join users u on u.id = a.user_id
			GROUP BY a.survey_id, a.id, a."date", a.ai_generated_status, u.id
			order by a."date" 
	) a on a.survey_id = s.id
	inner join landfills l on l.id = s.landfill_id 
	inner join regions r on r.id = l.region_id 
	inner join users u2 on u2.id = s.user_id
	inner  join(
		select 
		s2.id,
		array_agg(i."path") filter(where i.raw_status = true) as raw_images, 
		array_agg(i."path") filter(where i.raw_status = false) as ai_images
		from surveys s2 
		left join images i on i.survey_id = s2.id
		where s2.id = $1
		group by s2.id
	) i on i.id = s.id
	where s.id = $1
	GROUP BY s.id, r.id, u2.id, i.raw_images, i.ai_images
	`

	q := pg.Query{
		Name:     "survey.GetByID",
		QueryRaw: query,
	}

	var resultDB repentities.SurveyResponse
	err := rep.client.PG().GetContext(ctx, &resultDB, q, surveyID)
	if err != nil {
		return nil, err
	}
	if len(resultDB.Audits) > 0 {
		resultDB.ViolationsCount = resultDB.Audits[len(resultDB.Audits)-1].TrueViolationsCount
	}
	result := rep.converter.Survey.ToLogic(resultDB)

	if len(result.Audits) != 0 && result.Audits[0].ID == 0 {
		result.Audits = []logicentities.AuditResposeAll{}
	}
	if len(result.AiImages) != 0 && result.AiImages[0] == "" {
		result.AiImages = []string{}
	}
	if len(result.RawImages) != 0 && result.RawImages[0] == "" {
		result.RawImages = []string{}
	}

	return &result, nil
}

// DeleteSurvey ...
func (rep repository) DeleteSurvey(ctx context.Context, id uint64) error {
	s, err := rep.GetSurvey(ctx, id)
	if err != nil {
		return nil
	}
	for _, a := range s.Audits {
		err = rep.DeleteAudit(ctx, a.ID)
		if err != nil {
			return err
		}
	}
	// delete images
	err = rep.DeleteImagesBySurveyID(ctx, id)
	if err != nil {
		return err
	}

	deleteBuilder := squirrel.Delete("surveys").
		Where("id=?", id).PlaceholderFormat(sq.Dollar)

	// Generate the SQL and arguments
	query, v, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}
	q := pg.Query{
		Name:     "surveys.Delete",
		QueryRaw: query,
	}
	// Execute the delete operation
	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

func (rep repository) CheckUnique(ctx context.Context, date time.Time, landfillId uint64) bool {
	query := `
	select
	*
	from surveys
	where date = $1
	and landfill_id = $2
	`

	q := pg.Query{
		Name:     "survey.GetByID",
		QueryRaw: query,
	}

	var resultDB struct {
		ID         uint64    `json:"id"`
		LandfillID uint64    `json:"landfill_id"`
		Date       time.Time `json:"date"`
		UserID     uint64    `json:"user_id"`
	}
	err := rep.client.PG().GetContext(ctx, &resultDB, q, date, landfillId)
	if err != nil {
		return true
	}

	return resultDB.ID == 0
}
