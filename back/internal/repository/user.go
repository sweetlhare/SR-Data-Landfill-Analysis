package repository

import (
	"context"
	logicentities "svalka-service/internal/logic/entities"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

// CreateUser ...
func (rep repository) CreateUser(ctx context.Context, user logicentities.UserCreate) error {
	userDB := rep.converter.User.CreateToDB(user)
	columns, values := pg.ParseDbModel(userDB)
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "user.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// Login ...
func (rep repository) GetUserByEmail(ctx context.Context, email string) (user *logicentities.UserWithPass, err error) {
	// parse
	builder := squirrel.
		Select("*").
		From("users").
		Where("email = ?", email).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "users.GetByEmail",
		QueryRaw: query,
	}

	var respDB repentities.UserWithPass
	err = rep.client.PG().GetContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}

	resp := rep.converter.User.UserWithPassToLogic(respDB)

	return &resp, nil
}

// GetUser ...
func (rep repository) GetUser(ctx context.Context, userID uint64) (*logicentities.User, error) {
	if userID == 0 {
		return &logicentities.User{}, nil
	}
	// parse
	columns, _ := pg.ParseDbModel(repentities.User{})
	builder := squirrel.
		Select(columns...).
		From("users").
		Where("id = ?", userID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "users.Get",
		QueryRaw: query,
	}
	respDB := repentities.User{}
	err = rep.client.PG().GetContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}

	// convert
	user := rep.converter.User.ToLogic(respDB)

	return &user, nil
}
