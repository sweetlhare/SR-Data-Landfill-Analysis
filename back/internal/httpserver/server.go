package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HttpServer ...
type HttpServer interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	GetGinEngine() *gin.Engine
}

// httpServer ...
type httpServer struct {
	host            string
	shutdownTimeout time.Duration
	ginEngine       *gin.Engine
	server          *http.Server
}

// NewHttpServer ...
func NewHttpServer(_ context.Context, api HttpApi) (HttpServer, error) {
	httpServerConfig := getHttpServerConfig()

	// Set mode
	ginMode := GinModeFromString(httpServerConfig.Mode())
	gin.SetMode(ginMode)

	// Init gin engine
	ginEngine := gin.New()

	// Set recovery middleware
	ginEngine.Use(gin.Recovery())

	// Set status handler
	ginEngine.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})
	// Set API
	err := api.AddRoutes(ginEngine)
	if err != nil {
		return nil, err
	}
	// Set CORS
	setupCors(ginEngine)

	// Add static path for images
	ginEngine.Static("/images", "./images")

	// Init server
	server := &http.Server{
		Addr:    httpServerConfig.Host(),
		Handler: ginEngine,
	}

	return &httpServer{
		host:            httpServerConfig.Host(),
		shutdownTimeout: httpServerConfig.ShutdownTimeout(),
		ginEngine:       ginEngine,
		server:          server,
	}, nil
}

// Run ...
func (s *httpServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	if s.host != "" {
		// run
		wg.Add(1)
		go s.listenAndServe(ctx, wg)
		wg.Add(1)
		go s.softShutdown(ctx, wg)
	} else {
		log.Fatal(FailedStartServerError)
	}
}

// GetGinEngine ...
func (s *httpServer) GetGinEngine() *gin.Engine {
	return s.ginEngine
}

// listenAndServe ...
func (s *httpServer) listenAndServe(_ context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if s.server != nil {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(FailedListeningError, err)
		}
	}
}

// softShutdown ...
func (s *httpServer) softShutdown(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Server shutdown")
	// The context is used to inform the server it has s.shutdownTimeout to finish
	// the request it is currently handling
	ctxShutdown, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()
	err := s.close(ctxShutdown)
	if err != nil {
		log.Fatal(ServerForcedToShutdownError, zap.Error(err))
	}
	log.Print("Server exiting")
}

// close ...
func (s *httpServer) close(ctx context.Context) error {
	return nil
}
