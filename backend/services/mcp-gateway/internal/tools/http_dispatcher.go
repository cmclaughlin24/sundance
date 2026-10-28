package tools

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
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

func (d *HTTPDispatcher) Dispatch(ctx context.Context, tool *CompiledTool, args *ToolArgs) ([]byte, error) {
	url := args.replacePath(tool.BaseURL + tool.Path)

	body, err := args.serializeBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, tool.Method, url, body)
	if err != nil {
		return nil, err
	}

	req = args.appendQuery(req)
	req = setRequestHeaders(req, tool, args)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream request to %q failed: %w", tool.TargetName, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading upstream response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("upstream %s returned status code %d: %s", tool.TargetName, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func setRequestHeaders(req *http.Request, tool *CompiledTool, args *ToolArgs) *http.Request {
	req = tool.setHeaders(req)

	if len(args.Header) == 0 {
		return req
	}

	for k, v := range args.Header {
		if !tool.Headers.AllowHeader(k) {
			continue
		}

		req.Header.Set(k, fmt.Sprint(v))
	}

	return req
}
