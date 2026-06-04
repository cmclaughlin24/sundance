package strategies_test

import (
	"context"
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"sundance/backend/services/tenants/internal/core/strategies"
	"testing"
)

func TestDataLakeLookupStrategy_Lookup(t *testing.T) {
	dataLakeDS := &domain.DataSource{
		ID:   "ds-1",
		Type: domain.DataSourceTypeDataLake,
		Attributes: domain.DataLakeDataSourceAttributes{
			Query:        "SELECT value, label FROM affinity WHERE party = @party",
			RequiredKeys: []string{"party"},
			OptionalKeys: []string{"chapter"},
			Catalog:      "xenoblade",
			Schema:       "analytics",
			ValueField:   "value",
			LabelField:   "label",
			Limit:        100,
			TimeoutMs:    5000,
		},
	}

	tests := []struct {
		name    string
		ds      *domain.DataSource
		queryFn func(context.Context, domain.DataLakeDataSourceAttributes, map[string]any) ([]*domain.Lookup, error)
		want    []*domain.Lookup
		wantErr error
	}{
		{
			"should yield a list of lookups",
			dataLakeDS,
			func(_ context.Context, _ domain.DataLakeDataSourceAttributes, _ map[string]any) ([]*domain.Lookup, error) {
				return []*domain.Lookup{
					{Value: "shulk", Label: "Shulk"},
					{Value: "rex", Label: "Rex"},
				}, nil
			},
			[]*domain.Lookup{
				{Value: "shulk", Label: "Shulk"},
				{Value: "rex", Label: "Rex"},
			},
			nil,
		},
		{
			"should yield an empty list of lookups",
			dataLakeDS,
			func(_ context.Context, _ domain.DataLakeDataSourceAttributes, _ map[string]any) ([]*domain.Lookup, error) {
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
			dataLakeDS,
			func(_ context.Context, _ domain.DataLakeDataSourceAttributes, _ map[string]any) ([]*domain.Lookup, error) {
				return nil, errors.New("client error")
			},
			nil,
			errors.New("client error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := strategies.NewDataLakeLookupStrategy(logger, &ports.Clients{
				DataLake: &mockDataLakeClient{
					queryFn: func(ctx context.Context, attr domain.DataLakeDataSourceAttributes, params map[string]any) ([]*domain.Lookup, error) {
						if tt.queryFn != nil {
							return tt.queryFn(ctx, attr, params)
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
