package inmemory

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

func Bootstrap(logger *slog.Logger) *ports.Repository {
	return &ports.Repository{
		Database:    database.NewInMemoryDatabase(),
		Tenants:     newInMemoryTenantsRepository(logger),
		DataSources: newInMemoryDataSourceRepository(logger),
	}
}
