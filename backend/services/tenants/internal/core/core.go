package core

import (
	"context"
	"log/slog"

	"sundance/backend/pkg/cache"
	"sundance/backend/pkg/worker/elector"
	"sundance/backend/services/tenants/internal/core/ports"
)

type Cache interface {
	cache.CacheManager
	elector.CacheLocker
}

type Application struct {
	Logger     *slog.Logger
	API        *ports.API
	Cache      Cache
	repository *ports.Repository
}

func NewApplication(opts ...func(*Application)) *Application {
	var a Application
	for _, opt := range opts {
		opt(&a)
	}

	return &a
}

func WithCache(cache Cache) func(*Application) {
	return func(a *Application) {
		a.Cache = cache
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

func WithAPI(api *ports.API) func(*Application) {
	return func(a *Application) {
		a.API = api
	}
}

func (app *Application) Close(ctx context.Context) {
	if err := app.repository.Database.Close(ctx); err != nil {
		app.Logger.Error("an error occurred while closing the database connection", "error", err.Error())
	}
}
