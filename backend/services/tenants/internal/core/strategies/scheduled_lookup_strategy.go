package strategies

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type ScheduledLookupStrategy struct {
	logger *log.Logger
}

func NewScheduledLookupStrategy(logger *log.Logger) ports.LookupStrategy {
	return &ScheduledLookupStrategy{
		logger: logger,
	}
}

func (s *ScheduledLookupStrategy) Lookup(_ context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := getDataSourceAttributes[domain.ScheduledDataSourceAttributes](ds.Attributes)
	if err != nil {
		return nil, err
	}

	// TODO: Determine if making data lazy-loaded would make sense.
	return attr.Data, nil
}
