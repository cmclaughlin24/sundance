package dto_test

import (
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestDataSourceToResponse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		source *domain.DataSource
	}{
		{
			"should yield a response from a static data source",
			domain.HydrateDataSource(
				"ds-1", "f1-1", "Race Circuits",
				"The Monaco Grand Prix has been held since 1929 and is considered the crown jewel of Formula 1",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "monaco", Label: "Circuit de Monaco"},
						{Value: "silverstone", Label: "Silverstone Circuit"},
						{Value: "monza", Label: "Autodromo Nazionale Monza"},
					},
				},
				now, now.Add(24*time.Hour),
			),
		},
		{
			"should yield a response from a scheduled data source",
			domain.HydrateDataSource(
				"ds-2", "f1-1", "Driver Standings",
				"Michael Schumacher and Lewis Hamilton share the record of seven World Drivers Championships",
				domain.DataSourceTypeScheduled,
				domain.ScheduledDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "verstappen", Label: "Max Verstappen"},
						{Value: "hamilton", Label: "Lewis Hamilton"},
					},
					DataSourceHTTPRequest: domain.DataSourceHTTPRequest{
						URL:     "https://example.com/standings",
						Method:  "GET",
						Headers: map[string]string{"Authorization": "Bearer fia-token"},
					},
					IntervalHours:  24,
					ExpirationDate: now.Add(48 * time.Hour),
				},
				now, now.Add(24*time.Hour),
			),
		},
		{
			"should yield a response from a webhook data source",
			domain.HydrateDataSource(
				"ds-3", "f1-1", "Live Timing",
				"Formula 1 cars generate over 300 sensors worth of telemetry data during a race",
				domain.DataSourceTypeWebhook,
				domain.WebhookDataSourceAttributes{
					DataSourceHTTPRequest: domain.DataSourceHTTPRequest{
						URL:     "https://example.com/timing",
						Method:  "POST",
						Headers: map[string]string{"Authorization": "Bearer fia-token"},
					},
				},
				now, now.Add(24*time.Hour),
			),
		},
		{
			"should yield a response with empty optional fields",
			domain.HydrateDataSource(
				"ds-4", "f1-2", "Constructor Points",
				"",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{},
				},
				now, time.Time{},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := dto.DataSourceToResponse(tt.source)

			// Assert.
			if got == nil {
				t.Errorf("expected response but got nil")
				return
			}

			if got.ID != tt.source.ID {
				t.Errorf("expected ID %q but got %q", tt.source.ID, got.ID)
			}

			if got.TenantID != tt.source.TenantID {
				t.Errorf("expected TenantID %q but got %q", tt.source.TenantID, got.TenantID)
			}

			if got.Name != tt.source.Name {
				t.Errorf("expected name %q but got %q", tt.source.Name, got.Name)
			}

			if got.Description != tt.source.Description {
				t.Errorf("expected description %q but got %q", tt.source.Description, got.Description)
			}

			if got.Type != tt.source.Type {
				t.Errorf("expected type %q but got %q", tt.source.Type, got.Type)
			}

			if got.Attributes == nil {
				t.Errorf("expected non-nil Attributes")
			}

			if got.CreatedAt != tt.source.CreatedAt {
				t.Errorf("expected CreatedAt %v but got %v", tt.source.CreatedAt, got.CreatedAt)
			}

			if got.UpdatedAt != tt.source.UpdatedAt {
				t.Errorf("expected UpdatedAt %v but got %v", tt.source.UpdatedAt, got.UpdatedAt)
			}
		})
	}
}
