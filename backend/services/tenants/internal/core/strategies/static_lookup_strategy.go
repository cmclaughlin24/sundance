package strategies

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type StaticLookupStrategy struct {
	logger *log.Logger
}

func NewStaticLookupStrategy(logger *log.Logger) ports.LookupStrategy {
	return &StaticLookupStrategy{
		logger: logger,
	}
}

func (s *StaticLookupStrategy) Lookup(_ context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := getDataSourceAttributes[domain.StaticDataSourceAttributes](ds.Attributes)
	if err != nil {
		return nil, err
	}

	return attr.Data, nil
}
