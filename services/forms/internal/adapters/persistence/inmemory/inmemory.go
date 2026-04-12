package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	db := NewInMemoryDatabase()

	return &ports.Repository{
		Database: db,
		// Forms: NewInMemoryFormsRepository(db, logger),
	}
}
