package clients

import (
	"log/slog"
	"net/http"

	"sundance/backend/services/tenants/internal/core/ports"
)

type clientOptions struct {
	client *http.Client
	logger *slog.Logger
}

func Bootstrap(opts ...func(*clientOptions)) *ports.Clients {
	var co clientOptions
	for _, opt := range opts {
		opt(&co)
	}

	return &ports.Clients{
		Lookups: NewLookupClient(co.client, co.logger),
	}
}

func WithHTTPClient(client *http.Client) func(*clientOptions) {
	return func(co *clientOptions) {
		co.client = client
	}
}
func WithLogger(logger *slog.Logger) func(*clientOptions) {
	return func(co *clientOptions) {
		co.logger = logger
	}
}
