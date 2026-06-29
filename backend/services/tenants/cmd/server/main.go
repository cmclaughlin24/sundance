package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sundance/backend/pkg/cache"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/logger"
	"sundance/backend/pkg/worker"
	"sundance/backend/services/tenants/internal/adapters/clients"
	"sundance/backend/services/tenants/internal/adapters/persistence"
	"sundance/backend/services/tenants/internal/adapters/rest"
	"sundance/backend/services/tenants/internal/adapters/workers"
	"sundance/backend/services/tenants/internal/core"
	"sundance/backend/services/tenants/internal/core/services"
	"sundance/backend/services/tenants/internal/core/strategies"

	"github.com/caarlos0/env/v11"
)

type settings struct {
	Port        int                             `json:"port" env:"APP_PORT"`
	Persistence persistence.PersistenceSettings `json:"database" envPrefix:"APP_DATABASE_"`
	Cache       cache.CacheSettings             `json:"cache" envPrefix:"APP_CACHE_"`
	LogLevel    string                          `json:"logLevel" env:"APP_LOG_LEVEL"`
	Worker      workers.WorkerOptions           `json:"worker" envPrefix:"APP_WORKER_"`
	Server      rest.ServerOptions              `json:"server" envPrefix:"APP_SERVER_"`
}

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON/YAML file")
	flag.Parse()

	var settings settings
	if err := common.ReadSettings(*settingsPath, &settings); err != nil {
		slog.Warn("failed to read settings from file; defaulting to environment variables", "error", err)

		if err = env.Parse(&settings); err != nil {
			panic(err)
		}
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logger.LogLevelToLevel(settings.LogLevel),
	})
	l := slog.New(&worker.WorkerContextHandler{Handler: &logger.RequestContextHandler{Handler: handler}})

	r, err := persistence.Bootstrap(settings.Persistence, l)
	if err != nil {
		l.Error("failed to bootstrap persistance", "error", err.Error())
		panic(err)
	}

	cm, cacheClose, err := cache.Bootstrap(settings.Cache, l)
	if err != nil {
		l.Error("failed to bootstrap cache", "error", err.Error())
		panic(err)
	}
	defer cacheClose()

	c := clients.Bootstrap(clients.WithLogger(l), clients.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
	st := strategies.Bootstrap(strategies.WithLogger(l), strategies.WithClients(c))
	s := services.Bootstrap(services.WithLogger(l), services.WithRepository(r), services.WithStrategies(st), services.WithClients(c))
	app := core.NewApplication(core.WithLogger(l), core.WithRepository(r), core.WithAPI(s), core.WithCache(cm.(core.Cache)))

	defer app.Close(context.Background())
	mux := rest.NewRoutes(app, settings.Server)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	serverErrChan := make(chan error, 1)

	go func() {
		app.Logger.Info(fmt.Sprintf("application listening on :%d", settings.Port))

		if err := server.ListenAndServe(); err != nil {
			serverErrChan <- err
		}
	}()

	start, err := workers.Bootstrap(app, settings.Worker)
	if err != nil {
		panic(err)
	}
	wCtx, wCancel := context.WithCancel(context.Background())
	defer wCancel()
	start(wCtx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		app.Logger.Error(fmt.Sprintf("server error: %v", err))
	case sig := <-signalChan:
		app.Logger.Info(fmt.Sprintf("application received shutdown signal: %v", sig))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Info(fmt.Sprintf("application shutdown failed: %v", err))
		return
	}

	app.Logger.Info("application shutdown successful")
}
