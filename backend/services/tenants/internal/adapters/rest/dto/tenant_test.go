package dto_test

import (
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestTenantToResponse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		tenant *domain.Tenant
		want   *dto.TenantResponse
	}{
		{
			"should yield a TenantResponse",
			domain.HydrateTenant("f1-1", "Scuderia Ferrari", "Founded in 1929 by Enzo Ferrari, the oldest and most successful team in Formula 1 history", now, now.Add(24*time.Hour)),
			&dto.TenantResponse{
				ID:          "f1-1",
				Name:        "Scuderia Ferrari",
				Description: "Founded in 1929 by Enzo Ferrari, the oldest and most successful team in Formula 1 history",
				CreatedAt:   now,
				UpdatedAt:   now.Add(24 * time.Hour),
			},
		},
		{
			"should yield a TenantResponse with empty optional fields",
			domain.HydrateTenant("f1-2", "McLaren Racing", "", now, time.Time{}),
			&dto.TenantResponse{
				ID:          "f1-2",
				Name:        "McLaren Racing",
				Description: "",
				CreatedAt:   now,
				UpdatedAt:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := dto.TenantToResponse(tt.tenant)

			// Assert.
			if got == nil {
				t.Errorf("expected response but got nil")
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("expected ID %q but got %q", tt.want.ID, got.ID)
			}

			if got.Name != tt.want.Name {
				t.Errorf("expected name %q but got %q", tt.want.Name, got.Name)
			}

			if got.Description != tt.want.Description {
				t.Errorf("expected description %q but got %q", tt.want.Description, got.Description)
			}

			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("expected CreatedAt %v but got %v", tt.want.CreatedAt, got.CreatedAt)
			}

			if got.UpdatedAt != tt.want.UpdatedAt {
				t.Errorf("expected UpdatedAt %v but got %v", tt.want.UpdatedAt, got.UpdatedAt)
			}
		})
	}
}
