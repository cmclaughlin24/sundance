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
	fetchLookupsFn func(context.Context, string, string, map[string]string, map[string]any) ([]map[string]any, error)
}

func (c *mockLookupClient) FetchLookups(ctx context.Context, request domain.DataSourceHTTPRequest, params map[string]any) ([]map[string]any, error) {
	return c.fetchLookupsFn(ctx, request.Method, request.URL, request.Headers, params)
}

type mockDataLakeClient struct {
	queryFn func(context.Context, domain.DataLakeDataSourceAttributes, map[string]any) ([]*domain.Lookup, error)
}

func (c *mockDataLakeClient) Query(ctx context.Context, attr domain.DataLakeDataSourceAttributes, params map[string]any) ([]*domain.Lookup, error) {
	return c.queryFn(ctx, attr, params)
}
