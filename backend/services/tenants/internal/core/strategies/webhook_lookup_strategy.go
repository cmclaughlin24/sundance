package strategies

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type WebhookLookupStrategy struct {
	logger *slog.Logger
	client ports.HTTPClient
}

func NewWebhookLookupStrategy(logger *slog.Logger, client ports.HTTPClient) ports.LookupStrategy {
	return &WebhookLookupStrategy{
		logger: logger,
		client: client,
	}
}

func (s *WebhookLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := getDataSourceAttributes[domain.WebhookDataSourceAttributes](ds.Attributes)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, attr.Method, attr.URL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range attr.Headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var items []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	if err := httputil.DecodeJSONResponse(resp, &items); err != nil {
		return nil, err
	}

	lookups := make([]*domain.Lookup, 0, len(items))
	for _, item := range items {
		lookups = append(lookups, domain.NewLookup(item.Value, item.Label))
	}

	return lookups, nil
}
