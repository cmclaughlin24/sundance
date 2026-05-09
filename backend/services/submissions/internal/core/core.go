package core

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type Application struct {
	Logger     *slog.Logger
	Services   *ports.Services
	repository *ports.Repository
}

type applicationOptions struct {
	logger     *slog.Logger
	repository *ports.Repository
	services   *ports.Services
}

func NewApplication(opts ...func(*applicationOptions)) *Application {
	var ao applicationOptions
	for _, opt := range opts {
		opt(&ao)
	}

	return &Application{
		Logger:     ao.logger,
		Services:   ao.services,
		repository: ao.repository,
	}
}

func WithLogger(logger *slog.Logger) func(*applicationOptions) {
	return func(ao *applicationOptions) {
		ao.logger = logger
	}
}

func WithRepository(repository *ports.Repository) func(*applicationOptions) {
	return func(ao *applicationOptions) {
		ao.repository = repository
	}
}

func WithServices(services *ports.Services) func(*applicationOptions) {
	return func(ao *applicationOptions) {
		ao.services = services
	}
}

func (app *Application) Close() {
	if err := app.repository.Database.Close(); err != nil {
		app.Logger.Error("an error occurred while closing the database connection", "error", err.Error())
	}
}
