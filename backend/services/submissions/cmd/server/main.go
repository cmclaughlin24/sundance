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

	"github.com/cmclaughlin24/sundance/backend/pkg/cache"
	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/logger"
	"github.com/cmclaughlin24/sundance/backend/pkg/worker"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/persistence"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/rest"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/workers"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/services"
)

type settings struct {
	Port        int                             `json:"port"`
	Persistence persistence.PersistenceSettings `json:"persistence"`
	Cache       cache.CacheSettings             `json:"cache"`
	LogLevel    string                          `json:"log_level"`
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
	l := slog.New(&worker.WorkerContextHandler{Handler: &logger.RequestContextHandler{Handler: handler}})

	r, err := persistence.Bootstrap(settings.Persistence, l)
	if err != nil {
		l.Error("error", err.Error())
		panic(err)
	}

	cm, err := cache.Bootstrap(settings.Cache, l)
	if err != nil {
		l.Error("error", err.Error())
		panic(err)
	}

	s := services.Bootstrap(services.WithLogger(l), services.WithRepository(r))
	app := core.NewApplication(core.WithLogger(l), core.WithRepository(r), core.WithServices(s), core.WithCache(cm))

	defer app.Close(context.Background())
	mux := rest.NewRoutes(app)

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

	sw, err := workers.NewDataSourcesBackgroundWorker(app)
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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Info(fmt.Sprintf("application shutdown failed: %v", err))
		return
	}

	app.Logger.Info("application shutdown successful")
}
