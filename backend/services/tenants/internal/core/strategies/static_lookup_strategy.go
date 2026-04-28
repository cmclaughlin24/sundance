package strategies

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type StaticLookupStrategy struct{}

func NewStaticLookupStrategy() *StaticLookupStrategy {
	return &StaticLookupStrategy{}
}

func (s *StaticLookupStrategy) Lookup(_ context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	return nil, nil
}
