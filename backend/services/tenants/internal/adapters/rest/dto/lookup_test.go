package dto_test

import (
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
)

func TestLookupsToResponse(t *testing.T) {
	tests := []struct {
		name    string
		lookups []*domain.Lookup
		want    []*dto.LookupResponse
	}{
		{
			"should yield a list of lookup responses",
			[]*domain.Lookup{
				{Value: "monaco", Label: "Circuit de Monaco"},
				{Value: "silverstone", Label: "Silverstone Circuit"},
				{Value: "monza", Label: "Autodromo Nazionale Monza"},
			},
			[]*dto.LookupResponse{
				{Value: "monaco", Label: "Circuit de Monaco"},
				{Value: "silverstone", Label: "Silverstone Circuit"},
				{Value: "monza", Label: "Autodromo Nazionale Monza"},
			},
		},
		{
			"should yield an empty list",
			[]*domain.Lookup{},
			[]*dto.LookupResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := dto.LookupsToResponse(tt.lookups)

			// Assert.
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
