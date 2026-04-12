package services

import (
	"log"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

func Bootstrap(logger *log.Logger, repository *ports.Repository) *ports.Services {
	return &ports.Services{
		Tenants:     NewTenantsService(logger, repository),
		DataSources: NewDataSourcesService(logger, repository),
	}
}
