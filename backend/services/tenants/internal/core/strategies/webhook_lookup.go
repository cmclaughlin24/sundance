package strategies

import (
	"context"
	"fmt"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type WebhookLookupStrategy struct {
	baseLookupStrategy
	client ports.LookupClient
}

func NewWebhookLookupStrategy(logger *slog.Logger, clients *ports.Clients) ports.LookupStrategy {
	return &WebhookLookupStrategy{
		client: clients.Lookups,
		baseLookupStrategy: baseLookupStrategy{
			logger: logger,
		},
	}
}

func (s *WebhookLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource, params map[string]any) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.WebhookDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	if missing := s.missingRequiredKeys(attr.RequiredKeys, params); len(missing) > 0 {
		s.logger.ErrorContext(ctx, "webhook lookup missing required keys", "data_source_id", ds.ID, "missing_keys", missing)
		return nil, fmt.Errorf("%w: %v", domain.ErrMissingRequiredKeys, missing)
	}

	s.logger.DebugContext(ctx, "webhook lookup request", "data_source_id", ds.ID, "method", attr.Method, "url", attr.URL)

	rows, err := s.client.FetchLookups(ctx, attr.DataSourceHTTPRequest, params)
	if err != nil {
		return nil, err
	}

	lookups := s.toLookups(ctx, rows, attr.ValueField, attr.LabelField, ds.ID)

	s.logger.DebugContext(ctx, "webhook lookup resolved", "data_source_id", ds.ID, "count", len(lookups))

	return lookups, nil
}
