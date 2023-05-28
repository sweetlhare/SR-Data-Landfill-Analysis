package pgclient

import (
	"context"
	"svalka-service/pkg/closer"
	"svalka-service/pkg/pg"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	pgClient pg.Client
)

// GetPgClient ...
func GetPgClient(ctx context.Context) (pg.Client, error) {
	if pgClient == nil {
		pgCfg, err := pgxpool.ParseConfig(getPGConfig().DSN())
		if err != nil {
			return nil, DBConfigError
		}

		cl, err := pg.NewClient(ctx, pgCfg)
		if err != nil {
			return nil, PgClientError
		}

		err = cl.PG().Ping(ctx)
		if err != nil {
			return nil, PingError
		}
		closer.Add(cl.Close)

		pgClient = cl
	}

	return pgClient, nil
}
