package repository

import (
	"context"
	logicInterfaces "svalka-service/internal/logic/interfaces"
	pgClient "svalka-service/internal/pgclient"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"
)

type repository struct {
	client    pg.Client
	converter repentities.Converter
}

func NewRepository(ctx context.Context) (logicInterfaces.Repository, error) {
	postgresClient, err := pgClient.GetPgClient(ctx)
	if err != nil {
		return nil, err
	}
	return repository{
		client: postgresClient,
	}, nil
}
