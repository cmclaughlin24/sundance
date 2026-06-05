package clients_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sundance/backend/services/tenants/internal/adapters/clients"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
)

func TestLookupClient_FetchLookups(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		url     string
		headers map[string]string
		params  map[string]any
		doFn    func(*http.Request) (*http.Response, error)
		want    []map[string]any
		wantErr bool
	}{
		{
			"should yield a list of lookup rows",
			"GET",
			"https://example.com/blades",
			map[string]string{"Authorization": "Bearer driver-token"},
			nil,
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[{"value":"monado","label":"Monado"},{"value":"aegis","label":"Aegis"}]`)),
				}, nil
			},
			[]map[string]any{
				{"value": "monado", "label": "Monado"},
				{"value": "aegis", "label": "Aegis"},
			},
			false,
		},
		{
			"should yield an empty list",
			"GET",
			"https://example.com/blades",
			map[string]string{"Authorization": "Bearer driver-token"},
			nil,
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should yield an error when the request fails",
			"GET",
			"https://example.com/blades",
			nil,
			nil,
			func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("connection refused")
			},
			nil,
			true,
		},
		{
			"should yield an error when the response status is not ok",
			"GET",
			"https://example.com/blades",
			nil,
			nil,
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
			nil,
			true,
		},
		{
			"should yield an error when the response body is invalid",
			"GET",
			"https://example.com/blades",
			nil,
			nil,
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`not json`)),
				}, nil
			},
			nil,
			true,
		},
		{
			"should serialize params into the query string for GET",
			"GET",
			"https://example.com/blades",
			nil,
			map[string]any{"driver": "shulk", "season": 2},
			func(r *http.Request) (*http.Response, error) {
				q := r.URL.Query()
				if q.Get("driver") != "shulk" {
					return nil, errors.New("expected driver=shulk")
				}
				if q.Get("season") != "2" {
					return nil, errors.New("expected season=2")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should serialize params into the JSON body for POST",
			"POST",
			"https://example.com/blades",
			nil,
			map[string]any{"driver": "shulk"},
			func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("Content-Type") != "application/json" {
					return nil, errors.New("expected Content-Type application/json")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}
				var got map[string]any
				if err := json.Unmarshal(body, &got); err != nil {
					return nil, err
				}
				if got["driver"] != "shulk" {
					return nil, errors.New("expected driver=shulk in body")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should serialize params into the JSON body for PUT",
			"PUT",
			"https://example.com/blades",
			nil,
			map[string]any{"driver": "rex"},
			func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("Content-Type") != "application/json" {
					return nil, errors.New("expected Content-Type application/json")
				}
				body, _ := io.ReadAll(r.Body)
				if !strings.Contains(string(body), "rex") {
					return nil, errors.New("expected rex in body")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should serialize params into the JSON body for PATCH",
			"PATCH",
			"https://example.com/blades",
			nil,
			map[string]any{"driver": "noah"},
			func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("Content-Type") != "application/json" {
					return nil, errors.New("expected Content-Type application/json")
				}
				body, _ := io.ReadAll(r.Body)
				if !strings.Contains(string(body), "noah") {
					return nil, errors.New("expected noah in body")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should normalize lowercase method to uppercase",
			"get",
			"https://example.com/blades",
			nil,
			map[string]any{"driver": "shulk"},
			func(r *http.Request) (*http.Response, error) {
				if r.Method != "GET" {
					return nil, errors.New("expected method GET")
				}
				if r.URL.Query().Get("driver") != "shulk" {
					return nil, errors.New("expected driver in query string")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]map[string]any{},
			false,
		},
		{
			"should yield an error when params cannot be marshaled",
			"POST",
			"https://example.com/blades",
			nil,
			map[string]any{"bad": make(chan int)},
			func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("should not be called")
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			client := clients.NewLookupClient(&mockHttpClient{doFn: tt.doFn}, logger)

			// Act.
			got, gotErr := client.FetchLookups(context.Background(), domain.DataSourceHTTPRequest{Method: tt.method, URL: tt.url, Headers: tt.headers}, tt.params)

			// Assert.
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("expected %d rows but got %d", len(tt.want), len(got))
				return
			}

			for idx, want := range tt.want {
				for k, v := range want {
					if got[idx][k] != v {
						t.Errorf("row %d key %q: expected %v but got %v", idx, k, v, got[idx][k])
					}
				}
			}
		})
	}
}
