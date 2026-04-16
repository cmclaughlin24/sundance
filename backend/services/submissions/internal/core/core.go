package core

import (
	"log"
	"os"

	"github.com/cmclaughlin24/sundance/submissions/internal/adapters/persistence"
	"github.com/cmclaughlin24/sundance/submissions/internal/core/ports"
	"github.com/cmclaughlin24/sundance/submissions/internal/core/services"
)

type ApplicationSettings struct {
	Persistence persistence.PersistenceSettings `json:"persistence"`
}

type Application struct {
	Logger     *log.Logger
	Services   *ports.Services
	repository *ports.Repository
}

func NewApplication(settings *ApplicationSettings) (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	r, err := persistence.Bootstrap(settings.Persistence, logger)

	if err != nil {
		return nil, err
	}

	s := services.Bootstrap(logger, r)

	return &Application{
		Logger:     logger,
		Services:   s,
		repository: r,
	}, nil
}

func (app *Application) Close() {
	if err := app.repository.Database.Close(); err != nil {
		app.Logger.Fatalf("an error occurred while closing the database connection: %s", err.Error())
	}
}
