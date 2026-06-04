package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

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

func (c *LookupClient) FetchLookups(ctx context.Context, request domain.DataSourceRequest, params map[string]any) ([]map[string]any, error) {
	var body io.Reader
	method := strings.ToUpper(request.Method)

	if len(params) > 0 && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		b, err := c.setRequestBody(params)

		if err != nil {
			c.logger.ErrorContext(ctx, "lookup request body encode failed", "error", err)
			return nil, err
		}

		body = b
	}

	req, err := http.NewRequestWithContext(ctx, method, request.URL, body)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 && (method == http.MethodGet) {
		c.setQueryString(req, params)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "lookup request failed", "url", request.URL, "error", err)
		return nil, err
	}

	var data []map[string]any
	if err := httputil.DecodeJSONResponse(resp, &data); err != nil {
		c.logger.ErrorContext(ctx, "lookup request response decode failed", "error", err)
		return nil, err
	}

	return data, nil
}

func (c *LookupClient) setRequestBody(params map[string]any) (io.Reader, error) {
	jsonBytes, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBytes), nil
}

func (c *LookupClient) setQueryString(r *http.Request, params map[string]any) {
	query := r.URL.Query()

	for key, value := range params {
		query.Add(key, fmt.Sprint(value))
	}

	r.URL.RawQuery = query.Encode()
}
