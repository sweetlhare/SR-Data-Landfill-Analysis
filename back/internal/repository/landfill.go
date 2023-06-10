package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

// CreateLandfill ...
func (rep repository) CreateLandfill(ctx context.Context, landfill logicentities.LandfillCreate) (*logicentities.LandfillResponse, error) {
	landfillDB := rep.converter.Lanfill.ToDB(landfill)
	columns, values := pg.ParseDbModel(landfillDB)
	builder := sq.Insert("landfills").
		PlaceholderFormat(sq.Dollar).
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "landfill.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		rows.Close()
		return nil, err
	}

	rows.Next()
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}

	// get landfill
	return rep.GetLandfill(ctx, id)
}

// UpdateLandfill ...
func (rep repository) UpdateLandfill(ctx context.Context, landfill logicentities.LandfillUpdate) error {
	landfillDB := rep.converter.Lanfill.ToDB(landfill.LandfillCreate)
	columns, values := pg.ParseDbModel(landfillDB)
	builder := sq.Update("landfills")
	for i, c := range columns {
		builder = builder.Set(c, values[i])
	}
	builder = builder.PlaceholderFormat(sq.Dollar).Where("id=?", landfill.ID)
	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "landfill.Update",
		QueryRaw: query,
	}

	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

// GetLandfill ...
func (rep repository) GetLandfill(ctx context.Context, landfillID uint64) (*logicentities.LandfillResponse, error) {
	query := `
	select 
	l1.*, 
	l2.violations,
	l3.surveys
	from 
	(		SELECT 
			l2."id",
			l2."city",
			l2."name",
			l2."illegal",
			l2."address",
			l2."coordinates",
			l2."preview_image_path",
			l2."manager",
			l2."type",
			l2."cadastral_number",
			l2."illegal_cadastral_numbers",
			l2."cadastral_category",
			l2."usr_area",
			l2."area",
			jsonb_build_object('id', r.id, 'name', r."name") as region,
			coalesce (c.count, 0) as violations_count
		FROM landfills l2
			LEFT JOIN
			(
				SELECT l.id, count(*)
				FROM landfills l
				LEFT JOIN surveys s2 ON s2.landfill_id = l.id
				LEFT JOIN 
				(
					SELECT 
					s.id as s_id, 
					s."date", 
					s.landfill_id,
					a.id as a_id,
					a."date",
					ROW_NUMBER() OVER (
						PARTITION BY s.landfill_id 
						ORDER BY s."date" DESC, a."date" DESC
					) AS rn
					FROM surveys s
					INNER JOIN audits a ON a.survey_id = s.id
				) as s ON s.s_id = s2.id
				LEFT JOIN audits_to_violations atv ON s.a_id = atv.audit_id
				WHERE rn = 1 AND s.landfill_id = $1
				GROUP BY l.id, s2.id, atv.audit_id
			) as c ON l2.id = c.id
		INNER JOIN regions r ON r.id = l2.region_id
		WHERE l2.id = $1
	)
	l1 left join 
	(
	SELECT 
		l.id,
		json_agg(
			json_build_object
			(
				'id', v2.id,
				'title', v2.description,
				'default_status', v2.default_status,
				'status', v2.status
			)
		) as violations
	from landfills l
	LEFT JOIN(
		SELECT
			a.lanfill_id as lanfill_id,
			v.*, 
			v.id  = ANY(array_agg(atv.violation_id)) as status
		FROM (
			SELECT
				a.id,
				l.id as lanfill_id
			FROM audits a 
			INNER JOIN surveys s ON s.id = a.survey_id 
			INNER JOIN landfills l ON l.id  = s.landfill_id 
			WHERE l.id = $1
			ORDER BY s."date" DESC, a."date" DESC
			LIMIT 1
		) a 
		INNER JOIN audits_to_violations atv ON atv.audit_id = a.id
		JOIN violations v ON true
		GROUP BY a.lanfill_id, v.id 
		ORDER BY v.id
	) as v2 ON v2.lanfill_id = l.id
	WHERE l.id = $1
	group by l.id
	)
	l2 on l2.id = l1.id
	left join 
	(
		select 
		l."id",
		json_agg(jsonb_build_object('id', s.id, 'date', TO_CHAR(s."date", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'), 'landfill_id',s.landfill_id, 'author', s.author,'violations_count', s.violations_count)) as surveys  
		from landfills l
		left join 
		(
			select
			s.id, 
			s."date", 
			s.landfill_id,
			jsonb_build_object('id', u2.id, 'name', u2."name", 'position', u2."position",
			'role', u2."role", 'phone', u2.phone, 'email', u2.email) as author,
			count(*) as violations_count
			FROM public.surveys s
			left join 
			(
				SELECT 
				s.id as s_id,
				a.id as a_id,
				ROW_NUMBER() OVER 
				(
					PARTITION BY s.id
					ORDER BY a."date" desc
				) AS rn
				FROM surveys s
				inner JOIN audits a on a.survey_id = s.id
			) as s2 on s2.s_id = s.id
			left join audits_to_violations atv on s2.a_id = atv.audit_id
			inner join users u2 on u2.id = s.user_id 
			where rn = 1 and s.landfill_id = $1
			GROUP BY s.id, u2.id
			ORDER BY s."date"
		) s on s.landfill_id = l.id
		where l.id = $1
		group by l.id
	)
	l3 on l3.id = l1.id;
	`

	q := pg.Query{
		Name:     "landfills.GetByID",
		QueryRaw: query,
	}

	var resultDB repentities.LandfillResponse
	err := rep.client.PG().GetContext(ctx, &resultDB, q, landfillID)
	if err != nil {
		return nil, err
	}

	result := rep.converter.Lanfill.ToLogic(resultDB)

	if len(result.Surveys) != 0 && result.Surveys[0].ID == 0 {
		result.Surveys = []logicentities.SurveyResponseAll{}
	}
	if len(result.Violations) != 0 && result.Violations[0].Title == "" {
		result.Violations = []logicentities.Violation{}
	}
	return &result, nil
}

// GetAllLandfills ...
func (rep repository) GetAllLandfills(ctx context.Context, regionID uint64) ([]logicentities.LandfillAllResponse, error) {
	query := `
	select 
	l2."id",
	l2."city",
	l2."name",
	l2."illegal",
	l2."address",
	l2."coordinates",
	l2."preview_image_path",
	l2."type",
	l2."cadastral_number",
	l2."illegal_cadastral_numbers",
	l2."cadastral_category",
	l2."usr_area",
	l2."area",
	jsonb_build_object('id', r.id, 'name', r."name") as region,
	coalesce (c.count, 0) as violations_count
	from landfills l2
	left join
	(select l.id, count(*)
	from landfills l
	left join surveys s2 on s2.landfill_id = l.id
	left join 
	(
		SELECT 
		s.id as s_id, 
		s."date", 
		s.landfill_id,
		a.id as a_id,
		a."date",
		ROW_NUMBER() OVER (
			PARTITION BY s.landfill_id 
			ORDER BY s."date" desc, a."date" desc
		) AS rn
		FROM surveys s
		inner JOIN audits a on a.survey_id = s.id) as s on s.s_id = s2.id
		left join audits_to_violations atv on s.a_id = atv.audit_id
		where rn = 1
		GROUP BY l.id, s2.id, atv.audit_id
	) as c on l2.id = c.id
	inner join regions r on r.id = l2.region_id  
	where l2.region_id = $1
	`

	q := pg.Query{
		Name:     "landfills.GetByRegion",
		QueryRaw: query,
	}

	var resultDB []repentities.LandfillResponseAll
	err := rep.client.PG().SelectContext(ctx, &resultDB, q, regionID)
	if err != nil {
		return nil, err
	}

	result := rep.converter.Lanfill.AllToLogicArray(resultDB...)

	return result, nil
}

// DeleteLandfill ...
func (rep repository) DeleteLandfill(ctx context.Context, id uint64) error {
	l, err := rep.GetLandfill(ctx, id)
	if err != nil {
		return err
	}
	for _, s := range l.Surveys {
		err = rep.DeleteSurvey(ctx, s.ID)
		if err != nil {
			return err
		}
	}

	deleteBuilder := squirrel.Delete("landfills").
		Where("id=?", id).PlaceholderFormat(sq.Dollar)

	// Generate the SQL and arguments
	query, v, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}
	q := pg.Query{
		Name:     "landfill.Delete",
		QueryRaw: query,
	}
	// Execute the delete operation
	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}
