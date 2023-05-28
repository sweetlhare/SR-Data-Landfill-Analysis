package app

import (
	"context"
	aiClient "svalka-service/internal/aiclient"
	httpApi "svalka-service/internal/api/http"
	cdnClient "svalka-service/internal/cdnclient"
	"svalka-service/internal/config"
	"svalka-service/internal/httpserver"
	logic "svalka-service/internal/logic"
	repository "svalka-service/internal/repository"
	sessionrep "svalka-service/internal/sessionrep"
)

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		config.Init,
		a.initRepository,
		a.initSessionRepository,
		a.initCdnClient,
		a.initAiClient,
		a.initLogic,
		a.initHttpApi,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initPostgres ...
func (a *App) initRepository(ctx context.Context) (err error) {
	a.repository, err = repository.NewRepository(ctx)
	return err
}

// initSessionRepository ...
func (a *App) initSessionRepository(ctx context.Context) (err error) {
	a.sessionRepository, err = sessionrep.NewSessionRep(ctx)
	return err
}

// initCdnClient ...
func (a *App) initCdnClient(ctx context.Context) (err error) {
	a.cdnClient, err = cdnClient.NewCdnClient(ctx)
	return err
}

// initAiClient ...
func (a *App) initAiClient(ctx context.Context) (err error) {
	a.aiClient, err = aiClient.NewAiClient(ctx)
	return err
}

// initLogic ...
func (a *App) initLogic(ctx context.Context) (err error) {
	a.logic, err = logic.NewLogic(
		ctx,
		a.repository,
		a.sessionRepository,
		a.aiClient,
		a.cdnClient,
	)
	return err
}

// initHttpApi ...
func (a *App) initHttpApi(ctx context.Context) (err error) {
	a.httpApi, err = httpApi.NewHttpApi(ctx, a.logic)
	return err
}

// initHttpServer ...
func (a *App) initHttpServer(ctx context.Context) (err error) {
	a.httpServer, err = httpserver.NewHttpServer(ctx, a.httpApi)
	return err
}
