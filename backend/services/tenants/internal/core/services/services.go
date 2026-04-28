package services

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/strategies"
)

func Bootstrap(logger *log.Logger, repository *ports.Repository) *ports.Services {
	lookupStrategies := stratreg.New[domain.DataSourceType, LookupStrategy]().
		Set(domain.DataSourceTypeStatic, strategies.NewStaticLookupStrategy())

	return &ports.Services{
		Tenants:     NewTenantsService(logger, repository),
		DataSources: NewDataSourcesService(logger, repository, lookupStrategies),
	}
}
