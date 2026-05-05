package inmemory

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

func Bootstrap(logger *slog.Logger) *ports.Repository {
	return &ports.Repository{
		Database: database.NewInMemoryDatabase(),
		Forms:    NewInMemoryFormsRepository(logger),
		Versions: NewInMemoryVersionsRepository(logger),
	}
}
