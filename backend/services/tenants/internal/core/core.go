package core

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/cache"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type Application struct {
	Logger     *slog.Logger
	Services   *ports.Services
	Cache      cache.CacheManager
	repository *ports.Repository
}

func NewApplication(opts ...func(*Application)) *Application {
	var a Application
	for _, opt := range opts {
		opt(&a)
	}

	return &a
}

func WithCache(manager cache.CacheManager) func(*Application) {
	return func(a *Application) {
		a.Cache = manager
	}
}

func WithLogger(logger *slog.Logger) func(*Application) {
	return func(a *Application) {
		a.Logger = logger
	}
}

func WithRepository(repository *ports.Repository) func(*Application) {
	return func(a *Application) {
		a.repository = repository
	}
}

func WithServices(services *ports.Services) func(*Application) {
	return func(a *Application) {
		a.Services = services
	}
}

func (app *Application) Close(ctx context.Context) {
	if err := app.repository.Database.Close(ctx); err != nil {
		app.Logger.Error("an error occurred while closing the database connection", "error", err.Error())
	}
}
