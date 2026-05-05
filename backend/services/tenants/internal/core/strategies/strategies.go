package strategies

import (
	"errors"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

var (
	ErrDataSourceStrategyMismatch = errors.New("data source type and attributes mismatch; strategy cannot process")
)

func Bootstrap(logger *slog.Logger, _ *ports.Repository, client ports.HTTPClient) *ports.Strategies {
	lookupStrategies := stratreg.New[domain.DataSourceType, ports.LookupStrategy]().
		Set(domain.DataSourceTypeStatic, NewStaticLookupStrategy(logger)).
		Set(domain.DataSourceTypeScheduled, NewScheduledLookupStrategy(logger)).
		Set(domain.DataSourceTypeWebhook, NewWebhookLookupStrategy(logger, client))

	return &ports.Strategies{
		Lookups: lookupStrategies,
	}
}

func getDataSourceAttributes[T domain.DataSourceAttributes](attr domain.DataSourceAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrDataSourceStrategyMismatch
	}
}
