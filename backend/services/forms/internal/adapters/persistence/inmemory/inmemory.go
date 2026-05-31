package inmemory

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/ports"
)

func Bootstrap(logger *slog.Logger) *ports.Repository {
	return &ports.Repository{
		Database:      database.NewInMemoryDatabase(),
		Tags: newInMemoryTagRepository(logger),
		Forms:         newInMemoryFormsRepository(logger),
		FormVersions:  newInMemoryFormVersionsRepository(logger),
		Submissions:   newInMemorySubmissionsRepository(logger),
	}
}
