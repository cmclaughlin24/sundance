package core

import (
	"log/slog"
	"os"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type Application struct {
	Logger     *slog.Logger
	Services   *ports.Services
	repository *ports.Repository
}

func NewApplication(logger *slog.Logger, repository *ports.Repository, services *ports.Services) *Application {

	return &Application{
		Logger:     logger,
		Services:   services,
		repository: repository,
	}
}

func (app *Application) Close() {
	if err := app.repository.Database.Close(); err != nil {
		app.Logger.Error("an error occurred while closing the database connection", "error", err.Error())
		os.Exit(1)
	}
}
