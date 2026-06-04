package strategies_test

import (
	"context"
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"sundance/backend/services/tenants/internal/core/strategies"
	"testing"
)

func TestWebhookLookupStrategy_Lookup(t *testing.T) {
	webhookDS := &domain.DataSource{
		ID:   "ds-1",
		Type: domain.DataSourceTypeWebhook,
		Attributes: domain.WebhookDataSourceAttributes{
			URL:     "https://example.com/pokemon",
			Method:  "GET",
			Headers: map[string]string{"Authorization": "Bearer token"},
		},
	}

	tests := []struct {
		name           string
		ds             *domain.DataSource
		fetchLookupsFn func(context.Context, string, string, map[string]string) ([]*domain.Lookup, error)
		want           []*domain.Lookup
		wantErr        error
	}{
		{
			"should yield a list of lookups",
			webhookDS,
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return []*domain.Lookup{
					{Value: "bulbasaur", Label: "Bulbasaur"},
					{Value: "squirtle", Label: "Squirtle"},
				}, nil
			},
			[]*domain.Lookup{
				{Value: "bulbasaur", Label: "Bulbasaur"},
				{Value: "squirtle", Label: "Squirtle"},
			},
			nil,
		},
		{
			"should yield an empty list of lookups",
			webhookDS,
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return []*domain.Lookup{}, nil
			},
			[]*domain.Lookup{},
			nil,
		},
		{
			"should yield an error when the attribute type mismatches",
			&domain.DataSource{},
			nil,
			nil,
			domain.ErrDataSourceAttributeMismatch,
		},
		{
			"should yield an error when the client fails",
			webhookDS,
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return nil, errors.New("client error")
			},
			nil,
			errors.New("client error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := strategies.NewWebhookLookupStrategy(logger, &ports.Clients{
				Lookups: &mockLookupClient{
					fetchLookupsFn: func(ctx context.Context, method, url string, headers map[string]string) ([]*domain.Lookup, error) {
						if tt.fetchLookupsFn != nil {
							return tt.fetchLookupsFn(ctx, method, url, headers)
						}
						return nil, nil
					},
				},
			})

			// Act.
			got, gotErr := s.Lookup(context.Background(), tt.ds, nil)

			// Assert.
			if tt.wantErr != nil {
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
				t.Errorf("expected %d items but got %d", len(tt.want), len(got))
			}
		})
	}
}
