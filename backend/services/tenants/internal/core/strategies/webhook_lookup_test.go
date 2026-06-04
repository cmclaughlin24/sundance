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
			DataSourceRequest: domain.DataSourceRequest{
				URL:        "https://example.com/pokemon",
				Method:     "GET",
				Headers:    map[string]string{"Authorization": "Bearer token"},
				ValueField: "value",
				LabelField: "label",
			},
		},
	}

	customFieldsDS := &domain.DataSource{
		ID:   "ds-2",
		Type: domain.DataSourceTypeWebhook,
		Attributes: domain.WebhookDataSourceAttributes{
			DataSourceRequest: domain.DataSourceRequest{
				URL:        "https://example.com/pokemon",
				Method:     "GET",
				ValueField: "id",
				LabelField: "name",
			},
		},
	}

	requiredKeysDS := &domain.DataSource{
		ID:   "ds-3",
		Type: domain.DataSourceTypeWebhook,
		Attributes: domain.WebhookDataSourceAttributes{
			DataSourceRequest: domain.DataSourceRequest{
				URL:        "https://example.com/pokemon",
				Method:     "GET",
				ValueField: "value",
				LabelField: "label",
			},
			RequiredKeys: []string{"region"},
		},
	}

	tests := []struct {
		name           string
		ds             *domain.DataSource
		params         map[string]any
		fetchLookupsFn func(context.Context, string, string, map[string]string, map[string]any) ([]map[string]any, error)
		wantLen        int
		wantFirst      *domain.Lookup
		wantErr        error
	}{
		{
			"should yield a list of lookups",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "bulbasaur", "label": "Bulbasaur"},
					{"value": "squirtle", "label": "Squirtle"},
				}, nil
			},
			2,
			&domain.Lookup{Value: "bulbasaur", Label: "Bulbasaur"},
			nil,
		},
		{
			"should yield an empty list of lookups",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{}, nil
			},
			0,
			nil,
			nil,
		},
		{
			"should yield an error when the attribute type mismatches",
			&domain.DataSource{},
			nil,
			nil,
			0,
			nil,
			domain.ErrDataSourceAttributeMismatch,
		},
		{
			"should yield an error when the client fails",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return nil, errors.New("client error")
			},
			0,
			nil,
			errors.New("client error"),
		},
		{
			"should project rows using custom value and label fields",
			customFieldsDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"id": "001", "name": "Bulbasaur"},
				}, nil
			},
			1,
			&domain.Lookup{Value: "001", Label: "Bulbasaur"},
			nil,
		},
		{
			"should skip rows missing the value field",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "bulbasaur", "label": "Bulbasaur"},
					{"label": "No Value"},
				}, nil
			},
			1,
			&domain.Lookup{Value: "bulbasaur", Label: "Bulbasaur"},
			nil,
		},
		{
			"should skip rows missing the label field",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "bulbasaur", "label": "Bulbasaur"},
					{"value": "no-label"},
				}, nil
			},
			1,
			&domain.Lookup{Value: "bulbasaur", Label: "Bulbasaur"},
			nil,
		},
		{
			"should skip rows where the label is not a string",
			webhookDS,
			nil,
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "bulbasaur", "label": "Bulbasaur"},
					{"value": "weird", "label": 42},
				}, nil
			},
			1,
			&domain.Lookup{Value: "bulbasaur", Label: "Bulbasaur"},
			nil,
		},
		{
			"should yield an error when a required key is missing",
			requiredKeysDS,
			nil,
			nil,
			0,
			nil,
			domain.ErrMissingRequiredKeys,
		},
		{
			"should succeed when required keys are provided",
			requiredKeysDS,
			map[string]any{"region": "kanto"},
			func(_ context.Context, _, _ string, _ map[string]string, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "bulbasaur", "label": "Bulbasaur"},
				}, nil
			},
			1,
			&domain.Lookup{Value: "bulbasaur", Label: "Bulbasaur"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := strategies.NewWebhookLookupStrategy(logger, &ports.Clients{
				Lookups: &mockLookupClient{
					fetchLookupsFn: func(ctx context.Context, method, url string, headers map[string]string, params map[string]any) ([]map[string]any, error) {
						if tt.fetchLookupsFn != nil {
							return tt.fetchLookupsFn(ctx, method, url, headers, params)
						}
						return nil, nil
					},
				},
			})

			// Act.
			got, gotErr := s.Lookup(context.Background(), tt.ds, tt.params)

			// Assert.
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if errors.Is(tt.wantErr, domain.ErrDataSourceAttributeMismatch) ||
					errors.Is(tt.wantErr, domain.ErrMissingRequiredKeys) {
					if !errors.Is(gotErr, tt.wantErr) {
						t.Errorf("expected error wrapping %v but got %v", tt.wantErr, gotErr)
					}
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("expected %d lookups but got %d", tt.wantLen, len(got))
				return
			}

			if tt.wantFirst != nil {
				if got[0].Value != tt.wantFirst.Value {
					t.Errorf("expected first value %v but got %v", tt.wantFirst.Value, got[0].Value)
				}
				if got[0].Label != tt.wantFirst.Label {
					t.Errorf("expected first label %q but got %q", tt.wantFirst.Label, got[0].Label)
				}
			}
		})
	}
}
