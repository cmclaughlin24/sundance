package persistence

import (
	"fmt"
	"log"

	"github.com/cmclaughlin24/sundance/tenants/internal/adapters/persistence/inmemory"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

type PersistanceDriver string

const (
	PersistanceDriverInMemory PersistanceDriver = "in-memory"
)

type bootstrapFn func(PersistanceOptions, *log.Logger) (*ports.Repository, error)

type PersistanceOptions any

type PersistanceSettings struct {
	Driver  PersistanceDriver  `json:"driver"`
	Options PersistanceOptions `json:"options"`
}

func Bootstrap(settings PersistanceSettings, logger *log.Logger) (*ports.Repository, error) {
	var fn bootstrapFn

	switch settings.Driver {
	case PersistanceDriverInMemory:
		fn = bootstrapInMemory
	}

	if fn == nil {
		return nil, fmt.Errorf("unknown persistence driver: %s", settings.Driver)
	}

	return fn(settings.Options, logger)
}

func bootstrapInMemory(_ PersistanceOptions, logger *log.Logger) (*ports.Repository, error) {
	return inmemory.Bootstrap(logger), nil
}
