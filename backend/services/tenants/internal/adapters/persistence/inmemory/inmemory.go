package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	return &ports.Repository{
		Database:    database.NewInMemoryDatabase(),
		Tenants:     NewInmemoryTenantRepository(logger),
		DataSources: NewInmemoryDataSourceRepository(logger),
	}
}
