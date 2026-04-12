package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
)

type InMemoryFormsRepository struct {
	database ports.Database
	logger   *log.Logger
}

func NewInMemoryFormsRepository(database ports.Database, logger *log.Logger) *InMemoryFormsRepository {
	return &InMemoryFormsRepository{
		database: database,
		logger:   logger,
	}
}
