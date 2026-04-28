package core

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type Application struct {
	Logger     *log.Logger
	Services   *ports.Services
	repository *ports.Repository
}

func NewApplication(logger *log.Logger, repository *ports.Repository, services *ports.Services) *Application {

	return &Application{
		Logger:     logger,
		Services:   services,
		repository: repository,
	}
}

func (app *Application) Close() {
	if err := app.repository.Database.Close(); err != nil {
		app.Logger.Fatalf("an error occurred while closing the database connection: %s", err.Error())
	}
}
