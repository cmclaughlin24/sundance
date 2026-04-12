package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/common/database"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	return &ports.Repository{
		Database:    database.NewInMemoryDatabase(),
		Tenants:     NewInmemoryTenantRepository(logger),
		DataSources: NewInmemoryDataSourceRepository(logger),
	}
}
