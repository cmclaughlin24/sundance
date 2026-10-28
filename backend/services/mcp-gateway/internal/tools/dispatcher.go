package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Dispatcher interface {
	Dispatch(context.Context, *CompiledTool, *ToolArgs) ([]byte, error)
}

type ToolArgs struct {
	Header map[string]any `json:"header"`
	Path   map[string]any `json:"path"`
	Query  map[string]any `json:"query"`
	Body   map[string]any `json:"body"`
}

func (ta *ToolArgs) replacePath(url string) string {
	if ta.Path == nil || len(ta.Path) == 0 {
		return url
	}

	for k, v := range ta.Path {
		url = strings.ReplaceAll(url, "{"+k+"}", fmt.Sprint(v))
	}

	return url
}

func (ta *ToolArgs) appendQuery(req *http.Request) *http.Request {
	if ta.Query == nil || len(ta.Query) == 0 {
		return req
	}

	query := req.URL.Query()

	for k, v := range ta.Query {
		query.Add(k, fmt.Sprint(v))
	}

	req.URL.RawQuery = query.Encode()

	return req
}

func (ta *ToolArgs) setHeaders(req *http.Request) *http.Request {
	if len(ta.Header) == 0 {
		return req
	}

	for k, v := range ta.Header {
		req.Header.Set(k, fmt.Sprint(v))
	}

	return req
}

func (ta *ToolArgs) serializeBody() (io.Reader, error) {
	if ta.Body == nil || len(ta.Body) == 0 {
		return nil, nil
	}

	body, err := json.Marshal(ta.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewReader(body), nil
}
