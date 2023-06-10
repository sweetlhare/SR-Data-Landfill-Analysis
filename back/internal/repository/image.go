package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

// CreateImages ...
func (rep repository) CreateImages(ctx context.Context, images ...logicentities.ImageCreate) (err error) {
	columns, _ := pg.ParseDbModel(repentities.ImageCreate{})
	builder := squirrel.
		Insert("images").
		Columns(columns...).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range images {
		builder = builder.Values(i.SurveyID, i.Path, i.RawStatus)
	}

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "image.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	rows.Next()
	if rows.Err() != nil {
		return rows.Err()
	}
	defer rows.Close()

	return nil
}

// DeleteImagesBySurveyID ...
func (rep repository) DeleteImagesBySurveyID(ctx context.Context, id uint64) error {
	deleteBuilder := squirrel.Delete("images").
		Where("survey_id=?", id).PlaceholderFormat(sq.Dollar)

	// Generate the SQL and arguments
	query, v, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}
	q := pg.Query{
		Name:     "images.Delete",
		QueryRaw: query,
	}
	// Execute the delete operation
	_, err = rep.client.PG().ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}
