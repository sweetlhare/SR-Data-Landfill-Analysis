package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	sq "github.com/Masterminds/squirrel"
)

// CreateRegion ...
func (rep repository) CreateRegion(ctx context.Context, region logicentities.Region) error {
	regionDB := rep.converter.Region.CreateToDB(region)
	columns, values := pg.ParseDbModel(regionDB)
	builder := sq.Insert("regions").
		PlaceholderFormat(sq.Dollar).
		Columns(columns...).
		Values(values...)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "region.Create",
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

// GetAllRegions ...
func (rep repository) GetAllRegions(ctx context.Context) ([]logicentities.Region, error) {
	columns, _ := pg.ParseDbModel(repentities.Region{})
	builder := sq.Select(columns...).
		From("regions").
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "regions.GetAll",
		QueryRaw: query,
	}

	var resultDB []*repentities.Region
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}

	result := rep.converter.Region.ToLogicArray(resultDB...)
	return result, nil
}
