package clients_test

import (
	"context"
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
		doFn    func(*http.Request) (*http.Response, error)
		want    []*domain.Lookup
		wantErr bool
	}{
		{
			"should yield a list of lookups",
			"GET",
			"https://example.com/blades",
			map[string]string{"Authorization": "Bearer driver-token"},
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[{"value":"monado","label":"Monado"},{"value":"aegis","label":"Aegis"}]`)),
				}, nil
			},
			[]*domain.Lookup{
				{Value: "monado", Label: "Monado"},
				{Value: "aegis", Label: "Aegis"},
			},
			false,
		},
		{
			"should yield an empty list",
			"GET",
			"https://example.com/blades",
			map[string]string{"Authorization": "Bearer driver-token"},
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`[]`)),
				}, nil
			},
			[]*domain.Lookup{},
			false,
		},
		{
			"should yield an error when the request fails",
			"GET",
			"https://example.com/blades",
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
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`not json`)),
				}, nil
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
			got, gotErr := client.FetchLookups(context.Background(), tt.method, tt.url, tt.headers)

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
				t.Errorf("expected %d lookups but got %d", len(tt.want), len(got))
				return
			}

			for idx, want := range tt.want {
				if got[idx].Value != want.Value || got[idx].Label != want.Label {
					t.Errorf("expected %v but got %v", want, got[idx])
					break
				}
			}
		})
	}
}
