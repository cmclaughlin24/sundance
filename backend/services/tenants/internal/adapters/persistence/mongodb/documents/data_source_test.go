package documents_test

import (
	"sundance/backend/services/tenants/internal/adapters/persistence/mongodb/documents"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestToDataSourceDocument(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *domain.DataSource
	}{
		{
			"should yield a document from a static data source",
			domain.HydrateDataSource(
				"ds-1", "xc-1", "Monado Arts",
				"Techniques wielded by Shulk using the power of the Monado",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "speed", Label: "Speed"},
						{Value: "shield", Label: "Shield"},
						{Value: "buster", Label: "Buster"},
					},
				},
				now, now.Add(24*time.Hour),
			),
		},
		{
			"should yield a document from a scheduled data source",
			domain.HydrateDataSource(
				"ds-2", "xc-1", "Blade Resonance",
				"Blade awakening schedule from the core crystal registry",
				domain.DataSourceTypeScheduled,
				domain.ScheduledDataSourceAttributes{
					URL:           "https://example.com/blades",
					Method:        "GET",
					Headers:       map[string]string{"Authorization": "Bearer token"},
					IntervalHours: 12,
				},
				now, now.Add(24*time.Hour),
			),
		},
		{
			"should yield a document from a webhook data source",
			domain.HydrateDataSource(
				"ds-3", "xc-1", "Salvage Points",
				"Real-time salvage data from the Cloud Sea",
				domain.DataSourceTypeWebhook,
				domain.WebhookDataSourceAttributes{
					URL:     "https://example.com/salvage",
					Method:  "POST",
					Headers: map[string]string{"Authorization": "Bearer token"},
				},
				now, now.Add(24*time.Hour),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := documents.ToDataSourceDocument(tt.input)

			// Assert.
			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("expected document but got nil")
				return
			}

			if got.ID != string(tt.input.ID) {
				t.Errorf("expected ID %q but got %q", tt.input.ID, got.ID)
			}

			if got.TenantID != string(tt.input.TenantID) {
				t.Errorf("expected TenantID %q but got %q", tt.input.TenantID, got.TenantID)
			}

			if got.Name != tt.input.Name {
				t.Errorf("expected name %q but got %q", tt.input.Name, got.Name)
			}

			if got.Description != tt.input.Description {
				t.Errorf("expected description %q but got %q", tt.input.Description, got.Description)
			}

			if got.Type != string(tt.input.Type) {
				t.Errorf("expected type %q but got %q", tt.input.Type, got.Type)
			}

			if got.Attributes == nil {
				t.Errorf("expected non-nil Attributes")
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

func TestFromDataSourceDocument(t *testing.T) {
	now := time.Now()

	staticAttrRaw, _ := bson.Marshal(domain.StaticDataSourceAttributes{
		Data: []*domain.Lookup{
			{Value: "speed", Label: "Speed"},
			{Value: "shield", Label: "Shield"},
			{Value: "buster", Label: "Buster"},
		},
	})

	scheduledAttrRaw, _ := bson.Marshal(domain.ScheduledDataSourceAttributes{
		URL:           "https://example.com/blades",
		Method:        "GET",
		Headers:       map[string]string{"Authorization": "Bearer token"},
		IntervalHours: 12,
	})

	webhookAttrRaw, _ := bson.Marshal(domain.WebhookDataSourceAttributes{
		URL:     "https://example.com/salvage",
		Method:  "POST",
		Headers: map[string]string{"Authorization": "Bearer token"},
	})

	tests := []struct {
		name    string
		input   *documents.DataSourceDocument
		wantErr bool
	}{
		{
			"should yield a static data source from a document",
			&documents.DataSourceDocument{
				ID:          "ds-1",
				TenantID:    "xc-1",
				Name:        "Monado Arts",
				Description: "Techniques wielded by Shulk using the power of the Monado",
				Type:        "static",
				Attributes:  staticAttrRaw,
				CreatedAt:   now,
				UpdatedAt:   now.Add(24 * time.Hour),
			},
			false,
		},
		{
			"should yield a scheduled data source from a document",
			&documents.DataSourceDocument{
				ID:          "ds-2",
				TenantID:    "xc-1",
				Name:        "Blade Resonance",
				Description: "Blade awakening schedule from the core crystal registry",
				Type:        "scheduled",
				Attributes:  scheduledAttrRaw,
				CreatedAt:   now,
				UpdatedAt:   now.Add(24 * time.Hour),
			},
			false,
		},
		{
			"should yield a webhook data source from a document",
			&documents.DataSourceDocument{
				ID:          "ds-3",
				TenantID:    "xc-1",
				Name:        "Salvage Points",
				Description: "Real-time salvage data from the Cloud Sea",
				Type:        "webhook",
				Attributes:  webhookAttrRaw,
				CreatedAt:   now,
				UpdatedAt:   now.Add(24 * time.Hour),
			},
			false,
		},
		{
			"should yield an error for an unknown type",
			&documents.DataSourceDocument{
				ID:          "ds-4",
				TenantID:    "xc-1",
				Name:        "Nopon Commerce",
				Description: "Trade data from the Nopon Archsage",
				Type:        "unknown",
				Attributes:  staticAttrRaw,
				CreatedAt:   now,
				UpdatedAt:   now.Add(24 * time.Hour),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := documents.FromDataSourceDocument(tt.input)

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
				t.Errorf("expected data source but got nil")
				return
			}

			if string(got.ID) != tt.input.ID {
				t.Errorf("expected ID %q but got %q", tt.input.ID, got.ID)
			}

			if string(got.TenantID) != tt.input.TenantID {
				t.Errorf("expected TenantID %q but got %q", tt.input.TenantID, got.TenantID)
			}

			if got.Name != tt.input.Name {
				t.Errorf("expected name %q but got %q", tt.input.Name, got.Name)
			}

			if got.Description != tt.input.Description {
				t.Errorf("expected description %q but got %q", tt.input.Description, got.Description)
			}

			if string(got.Type) != tt.input.Type {
				t.Errorf("expected type %q but got %q", tt.input.Type, got.Type)
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
