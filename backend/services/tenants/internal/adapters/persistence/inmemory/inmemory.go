package inmemory

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/tenants/internal/core/ports"
)

func Bootstrap(logger *slog.Logger) *ports.Repository {
	return &ports.Repository{
		Database:    database.NewInMemoryDatabase(),
		Tenants:     newInMemoryTenantsRepository(logger),
		DataSources: newInMemoryDataSourceRepository(logger),
	}
}
