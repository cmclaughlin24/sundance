package strategies_test

import (
	"bytes"
	"context"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
)

var (
	buf    bytes.Buffer
	logger = slog.New(slog.NewTextHandler(&buf, nil))
)

type mockLookupClient struct {
	fetchLookupsFn func(context.Context, string, string, map[string]string) ([]*domain.Lookup, error)
}

func (c *mockLookupClient) FetchLookups(ctx context.Context, method, url string, headers map[string]string) ([]*domain.Lookup, error) {
	return c.fetchLookupsFn(ctx, method, url, headers)
}
