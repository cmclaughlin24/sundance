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

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/logger"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/persistence"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/rest"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/services"
)

type settings struct {
	Port        int                             `json:"port"`
	Persistence persistence.PersistenceSettings `json:"persistence"`
}

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON file")
	flag.Parse()

	var settings settings

	if err := common.ReadJSONFile(*settingsPath, &settings); err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, nil)
	l := slog.New(&logger.RequestContextHandler{Handler: handler})
	r, err := persistence.Bootstrap(settings.Persistence, l)

	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	s := services.Bootstrap(l, r)
	app := core.NewApplication(l, r, s)

	defer app.Close()
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
