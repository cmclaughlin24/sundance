package core

import (
	"log"
	"os"

	"github.com/cmclaughlin24/sundance/tenants/internal/adapters/persistence"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

type ApplicationSettings struct {
	Persistence persistence.PersistanceSettings `json:"persistence"`
}

type Application struct {
	Logger     *log.Logger
	Repository *ports.Repository
	Services   *ports.Services
}

func NewApplication(settings ApplicationSettings) (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_, err := persistence.Bootstrap(settings.Persistence, logger)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (app *Application) Close() {
	if err := app.Repository.Database.Close(); err != nil {
		log.Fatalf("an error occurred while closing the database connection: %s", err.Error())
	}
}
