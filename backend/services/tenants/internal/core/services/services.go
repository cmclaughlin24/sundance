package services

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

func Bootstrap(logger *log.Logger, repository *ports.Repository, strategies *ports.Strategies) *ports.Services {
	return &ports.Services{
		Tenants:     NewTenantsService(logger, repository),
		DataSources: NewDataSourcesService(logger, repository, strategies),
	}
}
