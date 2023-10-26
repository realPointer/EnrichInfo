package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/realPointer/EnrichInfo/config"
	v1 "github.com/realPointer/EnrichInfo/internal/controller/http/v1"
	"github.com/realPointer/EnrichInfo/internal/repo"
	"github.com/realPointer/EnrichInfo/internal/service"
	"github.com/realPointer/EnrichInfo/pkg/httpserver"
	"github.com/realPointer/EnrichInfo/pkg/logger"
	"github.com/realPointer/EnrichInfo/pkg/postgres"
)

func Run() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	l := logger.New(cfg.Log.Level)
	l.Info("Config and logger initialized")

	// Postgres
	l.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - pg.Pool.Ping: %w", err))
	}

	// Repositories
	l.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	l.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// HTTP Server
	l.Info("Initializing handlers and routes...")
	handler := chi.NewRouter()
	v1.NewRouter(handler, l, services)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	l.Info("Shutting down HTTP server...")
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
