package inmemory

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/ports"
)

func Bootstrap(logger *slog.Logger) *ports.Repository {
	outbox := newInMemoryOutbox(logger)

	return &ports.Repository{
		Database:     database.NewInMemoryDatabase(),
		Outbox:       outbox,
		Tags:         newInMemoryTagRepository(logger),
		TagVersions:  newInMemoryTagVersionsRepository(logger),
		Forms:        newInMemoryFormsRepository(logger),
		FormVersions: newInMemoryFormVersionsRepository(logger),
		Submissions:  newInMemorySubmissionsRepository(logger, outbox),
	}
}
