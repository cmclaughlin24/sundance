package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	return &ports.Repository{
		Database: database.NewInMemoryDatabase(),
		Forms:    NewInMemoryFormsRepository(logger),
	}
}
