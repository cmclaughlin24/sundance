package persistence

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/persistence/mongodb"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/persistence/inmemory"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type PersistenceDriver string

const (
	PersistenceDriverInMemory PersistenceDriver = "in-memory"
	PersistenceDriverMongodb  PersistenceDriver = "mongodb"
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
	case PersistenceDriverMongodb:
		fn = bootstrapMongoDB
	}

	if fn == nil {
		return nil, fmt.Errorf("unknown persistence driver: %s", settings.Driver)
	}

	return fn(settings.Options, logger)
}

func bootstrapInMemory(_ PersistenceOptions, logger *log.Logger) (*ports.Repository, error) {
	return inmemory.Bootstrap(logger), nil
}

func bootstrapMongoDB(o PersistenceOptions, logger *log.Logger) (*ports.Repository, error) {
	options, err := parseOptions[mongodb.MongoDBOpts](o)

	if err != nil {
		return nil, err
	}

	client, err := mongodb.Connect(
		mongodb.WithHost(options.Host),
		mongodb.WithPort(options.Port),
		mongodb.WithUsername(options.Username),
		mongodb.WithPassword(options.Password),
	)

	if err != nil {
		return nil, err
	}

	return mongodb.Bootstrap(client, logger), nil
}

func parseOptions[T PersistenceOptions](options PersistenceOptions) (T, error) {
	data, err := json.Marshal(options)

	if err != nil {
		return *new(T), err
	}

	var opts T

	if err := json.Unmarshal(data, &opts); err != nil {
		return *new(T), err
	}

	return opts, nil
}
