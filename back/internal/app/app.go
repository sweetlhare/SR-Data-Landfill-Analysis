package app

import (
	"context"
	"errors"
	"sync"

	"svalka-service/internal/httpserver"
	"svalka-service/pkg/closer"

	logicInterfaces "svalka-service/internal/logic/interfaces"
)

type App struct {
	httpServer        httpserver.HttpServer
	httpApi           httpserver.HttpApi
	aiClient          logicInterfaces.AiClient
	cdnClient         logicInterfaces.CdnClient
	logic             logicInterfaces.Logic
	repository        logicInterfaces.Repository
	sessionRepository logicInterfaces.SessionRepository
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	// run http server
	if a.httpServer == nil {
		return errors.New("http server not initialized")
	}
	wg := &sync.WaitGroup{}
	a.httpServer.Run(ctx, wg)
	wg.Wait()

	return nil
}
