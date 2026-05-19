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
	"sundance/backend/services/forms/internal/adapters/evaluators"
	"sundance/backend/services/forms/internal/adapters/persistence"
	"sundance/backend/services/forms/internal/adapters/rest"
	"sundance/backend/services/forms/internal/adapters/workers"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/services"
	"sundance/backend/services/forms/internal/core/strategies"
)

type settings struct {
	Port        int                             `json:"port"`
	Persistence persistence.PersistenceSettings `json:"persistence"`
	Cache       cache.CacheSettings             `json:"cache"`
	LogLevel    string                          `json:"log_level"`
	Host        string                          `json:"host"`
}

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON file")
	flag.Parse()

	var settings settings

	if err := common.ReadJSONFile(*settingsPath, &settings); err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logger.LogLevelToLevel(settings.LogLevel),
	})
	l := slog.New(&logger.RequestContextHandler{Handler: handler})

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

	st := strategies.Bootstrap(strategies.WithLogger(l))
	s := services.Bootstrap(services.WithLogger(l), services.WithRepository(r), services.WithStrategies(st), services.WithRuleEvaluator(&evaluators.ExprRuleEvaluator{}))
	app := core.NewApplication(core.WithLogger(l), core.WithRepository(r), core.WithServices(s), core.WithCache(cm.(core.Cache)))

	defer app.Close(context.Background())
	mux := rest.NewRoutes(app, settings.Host)

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

	sw, err := workers.NewSubmissionsBackgroundWorker(app)
	if err != nil {
		panic(err)
	}
	swCtx, swCancel := context.WithCancel(context.Background())
	defer swCancel()
	go sw.Start(swCtx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		app.Logger.Error(fmt.Sprintf("server error: %v", err))
	case sig := <-signalChan:
		app.Logger.Info(fmt.Sprintf("application received shutdown signal: %v", sig))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Info(fmt.Sprintf("application shutdown failed: %v", err))
		return
	}

	app.Logger.Info("application shutdown successful")
}
