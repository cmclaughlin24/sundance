package clients

import (
	"context"
	"log/slog"
	"net/http"

	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type LookupClient struct {
	client httpClient
	logger *slog.Logger
}

func NewLookupClient(client httpClient, logger *slog.Logger) ports.LookupClient {
	return &LookupClient{
		client: client,
		logger: logger,
	}
}

func (c *LookupClient) FetchLookups(ctx context.Context, method, url string, headers map[string]string) ([]*domain.Lookup, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "lookup request failed", "url", url, "error", err)
		return nil, err
	}

	var items []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	if err := httputil.DecodeJSONResponse(resp, &items); err != nil {
		c.logger.ErrorContext(ctx, "lookup request response decode failed", "error", err)
		return nil, err
	}

	lookups := make([]*domain.Lookup, 0, len(items))
	for _, item := range items {
		lookups = append(lookups, domain.NewLookup(item.Value, item.Label))
	}

	return lookups, nil
}
