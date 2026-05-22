package documents_test

import (
	"sundance/backend/services/tenants/internal/adapters/persistence/mongodb/documents"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestToTenantDocument(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *domain.Tenant
	}{
		{
			"should yield a TenantDocument",
			domain.HydrateTenant("xc-1", "Colony 9", "A defensive colony located near the foot of the Bionis", now, now.Add(24*time.Hour)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := documents.ToTenantDocument(tt.input)

			// Assert.
			if got == nil {
				t.Errorf("expected document but got nil")
				return
			}

			if got.ID != string(tt.input.ID) {
				t.Errorf("expected ID %q but got %q", tt.input.ID, got.ID)
			}

			if got.Name != tt.input.Name {
				t.Errorf("expected name %q but got %q", tt.input.Name, got.Name)
			}

			if got.Description != tt.input.Description {
				t.Errorf("expected description %q but got %q", tt.input.Description, got.Description)
			}

			if got.CreatedAt != tt.input.CreatedAt {
				t.Errorf("expected CreatedAt %v but got %v", tt.input.CreatedAt, got.CreatedAt)
			}

			if got.UpdatedAt != tt.input.UpdatedAt {
				t.Errorf("expected UpdatedAt %v but got %v", tt.input.UpdatedAt, got.UpdatedAt)
			}
		})
	}
}

func TestFromTenantDocument(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *documents.TenantDocument
	}{
		{
			"should yield a Tenant",
			&documents.TenantDocument{
				ID:          "xc-2",
				Name:        "Colony 6",
				Description: "A small colony situated on the Bionis' leg",
				CreatedAt:   now,
				UpdatedAt:   now.Add(48 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := documents.FromTenantDocument(tt.input)

			// Assert.
			if got == nil {
				t.Errorf("expected tenant but got nil")
				return
			}

			if string(got.ID) != tt.input.ID {
				t.Errorf("expected ID %q but got %q", tt.input.ID, got.ID)
			}

			if got.Name != tt.input.Name {
				t.Errorf("expected name %q but got %q", tt.input.Name, got.Name)
			}

			if got.Description != tt.input.Description {
				t.Errorf("expected description %q but got %q", tt.input.Description, got.Description)
			}

			if got.CreatedAt != tt.input.CreatedAt {
				t.Errorf("expected CreatedAt %v but got %v", tt.input.CreatedAt, got.CreatedAt)
			}

			if got.UpdatedAt != tt.input.UpdatedAt {
				t.Errorf("expected UpdatedAt %v but got %v", tt.input.UpdatedAt, got.UpdatedAt)
			}
		})
	}
}
