package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type WebhookLookupStrategy struct {
	logger *slog.Logger
	client ports.LookupClient
}

func NewWebhookLookupStrategy(logger *slog.Logger, clients *ports.Clients) ports.LookupStrategy {
	return &WebhookLookupStrategy{
		logger: logger,
		client: clients.Lookups,
	}
}

func (s *WebhookLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.WebhookDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	s.logger.DebugContext(ctx, "webhook lookup request", "data_source_id", ds.ID, "method", attr.Method, "url", attr.URL)

	lookups, err := s.client.FetchLookups(ctx, attr.Method, attr.URL, attr.Headers)
	if err != nil {
		return nil, err
	}

	s.logger.DebugContext(ctx, "webhook lookup resolved", "data_source_id", ds.ID, "count", len(lookups))

	return lookups, nil
}
