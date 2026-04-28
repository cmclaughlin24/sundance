package core

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/services"
)

type Application struct {
	Logger     *log.Logger
	Services   *ports.Services
	repository *ports.Repository
}

func NewApplication(logger *log.Logger, repository *ports.Repository) *Application {
	s := services.Bootstrap(logger, repository)

	return &Application{
		Logger:     logger,
		Services:   s,
		repository: repository,
	}
}

func (app *Application) Close() {
	if err := app.repository.Database.Close(); err != nil {
		app.Logger.Fatalf("an error occurred while closing the database connection: %s", err.Error())
	}
}
