package core

import (
	"log"

	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
)

type ApplicationSettings struct{}

type Application struct {
	Logger     *log.Logger
	Services   *ports.Services
	repository *ports.Repository
}

func NewApplication(settings *ApplicationSettings) (*Application, error) {
	return nil, nil
}
