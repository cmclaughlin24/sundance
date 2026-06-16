package persistence

import (
	"errors"
	"fmt"
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/tenants/internal/adapters/persistence/inmemory"
	"sundance/backend/services/tenants/internal/adapters/persistence/mongodb"
	"sundance/backend/services/tenants/internal/core/ports"
)

type PersistenceDriver string

const (
	PersistenceDriverInMemory PersistenceDriver = "in-memory"
	PersistenceDriverMongodb  PersistenceDriver = "mongodb"
)

type PersistenceSettings struct {
	Driver  PersistenceDriver     `json:"driver" env:"DRIVER"`
	MongoDB *database.MongoDBOpts `json:"mongodb,omitempty" envPrefix:"MONGODB_" env:",init"`
}

func Bootstrap(settings PersistenceSettings, logger *slog.Logger) (*ports.Repository, error) {
	switch settings.Driver {
	case PersistenceDriverInMemory:
		return bootstrapInMemory(logger)
	case PersistenceDriverMongodb:
		return bootstrapMongoDB(settings.MongoDB, logger)
	default:
		return nil, fmt.Errorf("unknown persistence driver: %s", settings.Driver)
	}
}

func bootstrapInMemory(logger *slog.Logger) (*ports.Repository, error) {
	return inmemory.Bootstrap(logger), nil
}

func bootstrapMongoDB(options *database.MongoDBOpts, logger *slog.Logger) (*ports.Repository, error) {
	if options == nil {
		return nil, errors.New("mongodb options are required for mongodb persistence driver")
	}

	client, err := database.ConnectMongoDB(
		database.MongoDBWithURI(options.URI),
	)

	if err != nil {
		return nil, err
	}

	return mongodb.Bootstrap(client, logger, options.DatabaseName), nil
}
