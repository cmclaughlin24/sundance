package persistence

import (
	"fmt"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/persistence/inmemory"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type PersistenceDriver string

const (
	PersistenceDriverInMemory PersistenceDriver = "in-memory"
)

type bootstrapFn func(PersistenceOptions, *log.Logger) (*ports.Repository, error)

type PersistenceOptions any

type PersistenceSettings struct {
	Driver  PersistenceDriver  `json:"driver"`
	Options PersistenceOptions `json:"options"`
}

func Bootstrap(settings PersistenceSettings, logger *log.Logger) (*ports.Repository, error) {
	var fn bootstrapFn

	switch settings.Driver {
	case PersistenceDriverInMemory:
		fn = bootstrapInMemory
	}

	if fn == nil {
		return nil, fmt.Errorf("unknown persistence driver: %s", settings.Driver)
	}

	return fn(settings.Options, logger)
}

func bootstrapInMemory(_ PersistenceOptions, logger *log.Logger) (*ports.Repository, error) {
	return inmemory.Bootstrap(logger), nil
}
