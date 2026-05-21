package strategies_test

import (
	"context"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/strategies"
	"testing"
)

func TestStaticLookupStrategy_Lookup(t *testing.T) {
	tests := []struct {
		name    string
		ds      *domain.DataSource
		wantErr error
	}{
		{
			"should yield a list of lookups",
			&domain.DataSource{
				Type: domain.DataSourceTypeStatic,
				Attributes: domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "pikachu", Label: "Pikachu"},
						{Value: "charizard", Label: "Charizard"},
					},
				},
			},
			nil,
		},
		{
			"should yield an empty list of lookups",
			&domain.DataSource{
				Type: domain.DataSourceTypeStatic,
				Attributes: domain.StaticDataSourceAttributes{
					Data: make([]*domain.Lookup, 0),
				},
			},
			nil,
		},
		{
			"should yield an error when the attribute type mismatches",
			&domain.DataSource{},
			domain.ErrDataSourceAttributeMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := strategies.NewStaticLookupStrategy(logger)

			// Act.
			got, gotErr := s.Lookup(context.Background(), tt.ds)

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

			attr, _ := tt.ds.Attributes.(domain.StaticDataSourceAttributes)

			if len(got) != len(attr.Data) {
				t.Errorf("expected %d items but got %d", len(attr.Data), len(got))
			}
		})
	}
}
