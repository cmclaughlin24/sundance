package clients_test

import (
	"bytes"
	"log/slog"
	"net/http"
)

var (
	buf    bytes.Buffer
	logger = slog.New(slog.NewTextHandler(&buf, nil))
)

type mockHttpClient struct {
	doFn func(*http.Request) (*http.Response, error)
}

func (c *mockHttpClient) Do(request *http.Request) (*http.Response, error) {
	return c.doFn(request)
}
