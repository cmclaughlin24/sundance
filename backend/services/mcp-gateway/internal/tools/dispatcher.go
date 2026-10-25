package tools

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Args struct {
	Path  map[string]any `json:"path"`
	Query map[string]any `json:"query"`
	Body  map[string]any `json:"body"`
}

type Dispatcher interface {
	Dispatch(context.Context, *CompiledTool, *Args) ([]byte, error)
}

type HTTPDispatcher struct {
	client httpClient
	logger *slog.Logger
}

func NewHTTPDispatcher(client httpClient, logger *slog.Logger) *HTTPDispatcher {
	return &HTTPDispatcher{
		client: client,
		logger: logger,
	}
}

func (d *HTTPDispatcher) Dispatch(ctx context.Context, tool *CompiledTool, args *Args) ([]byte, error) {
	return nil, errors.New("method not implemented")
}
