package domain_test

import (
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestNewTenant(t *testing.T) {
	tests := []struct {
		name        string
		tenantName  string
		description string
		wantErr     bool
	}{
		{
			"should create a tenant",
			"Kanto Pokemon League",
			"The original 151 Pokemon from the Kanto region",
			false,
		},
		{
			"should create a tenant with an empty description",
			"Johto Pokemon League",
			"",
			false,
		},
		{
			"should yield an error when the name is empty",
			"",
			"Missing name",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := domain.NewTenant(tt.tenantName, tt.description)

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

			if got == nil {
				t.Errorf("expected tenant but got nil")
				return
			}

			if got.ID == "" {
				t.Errorf("expected non-empty ID")
			}

			if got.Name != tt.tenantName {
				t.Errorf("expected name %q but got %q", tt.tenantName, got.Name)
			}

			if got.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, got.Description)
			}

			if got.CreatedAt.IsZero() {
				t.Errorf("expected non-zero CreatedAt")
			}
		})
	}
}

func TestHydrateTenant(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		id          domain.TenantID
		tenantName  string
		description string
		createdAt   time.Time
		updatedAt   time.Time
	}{
		{
			"should hydrate a tenant",
			"pkmn-1",
			"Hoenn Pokemon League",
			"Home to 135 new Pokemon introduced in Generation III",
			now,
			now.Add(24 * time.Hour),
		},
		{
			"should hydrate a tenant with empty optional fields",
			"pkmn-2",
			"Sinnoh Pokemon League",
			"",
			now,
			time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := domain.HydrateTenant(tt.id, tt.tenantName, tt.description, tt.createdAt, tt.updatedAt)

			// Assert.
			if got == nil {
				t.Errorf("expected tenant but got nil")
				return
			}

			if got.ID != tt.id {
				t.Errorf("expected ID %q but got %q", tt.id, got.ID)
			}

			if got.Name != tt.tenantName {
				t.Errorf("expected name %q but got %q", tt.tenantName, got.Name)
			}

			if got.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, got.Description)
			}

			if got.CreatedAt != tt.createdAt {
				t.Errorf("expected CreatedAt %v but got %v", tt.createdAt, got.CreatedAt)
			}

			if got.UpdatedAt != tt.updatedAt {
				t.Errorf("expected UpdatedAt %v but got %v", tt.updatedAt, got.UpdatedAt)
			}
		})
	}
}

func TestTenant_Update(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		tenant      *domain.Tenant
		tenantName  string
		description string
		wantErr     bool
	}{
		{
			"should update a tenant",
			domain.HydrateTenant("pkmn-1", "Unova Pokemon League", "old desc", now, time.Time{}),
			"Unova Champions League",
			"Introduced 156 new Pokemon in Generation V",
			false,
		},
		{
			"should update a tenant with an empty description",
			domain.HydrateTenant("pkmn-2", "Kalos Pokemon League", "old desc", now, time.Time{}),
			"Kalos Champions League",
			"",
			false,
		},
		{
			"should yield an error when the name is empty",
			domain.HydrateTenant("pkmn-3", "Alola Pokemon League", "old desc", now, time.Time{}),
			"",
			"some description",
			true,
		},
		{
			"should yield an error when the tenant is nil",
			nil,
			"Galar Pokemon League",
			"Generation VIII introduced the Wild Area",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			var originalName string
			if tt.tenant != nil {
				originalName = tt.tenant.Name
			}

			// Act.
			gotErr := tt.tenant.Update(tt.tenantName, tt.description)

			// Assert.
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
				}

				if tt.tenant == nil {
					if !errors.Is(gotErr, domain.ErrInvalidTenant) {
						t.Errorf("expected ErrInvalidTenant but got %v", gotErr)
					}
					return
				}

				if tt.tenant.Name != originalName {
					t.Errorf("expected name to remain %q but got %q", originalName, tt.tenant.Name)
				}

				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if tt.tenant.Name != tt.tenantName {
				t.Errorf("expected name %q but got %q", tt.tenantName, tt.tenant.Name)
			}

			if tt.tenant.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, tt.tenant.Description)
			}

			if tt.tenant.UpdatedAt.IsZero() {
				t.Errorf("expected non-zero UpdatedAt")
			}
		})
	}
}
