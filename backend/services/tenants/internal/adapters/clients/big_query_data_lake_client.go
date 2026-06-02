package clients

import (
	"context"
	"errors"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

var ErrBigQueryDataLakeNotConfigured = errors.New("big query data lake client not configured")

type BigQueryDataLakeClient struct {
	logger *slog.Logger
}

func NewBigQueryDataLakeClient(logger *slog.Logger) ports.DataLakeClient {
	return &BigQueryDataLakeClient{
		logger: logger,
	}
}

func (c *BigQueryDataLakeClient) Query(ctx context.Context, attr domain.DataLakeDataSourceAttributes, params map[string]any) ([]*domain.Lookup, error) {
	return nil, ErrBigQueryDataLakeNotConfigured
}
