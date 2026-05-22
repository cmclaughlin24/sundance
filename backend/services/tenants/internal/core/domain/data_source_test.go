package domain_test

import (
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestNewDataSource(t *testing.T) {
	tests := []struct {
		name           string
		tenantID       domain.TenantID
		dataSourceName string
		description    string
		dataSourceType domain.DataSourceType
		attributes     domain.DataSourceAttributes
		wantErr        error
	}{
		{
			"should create a static data source",
			"tenant-1",
			"Kanto Pokedex",
			"The original 151 Pokemon from the Kanto region",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "pikachu", Label: "Pikachu"},
					{Value: "charizard", Label: "Charizard"},
				},
			},
			nil,
		},
		{
			"should create a schedule data source",
			"tenant-1",
			"Unova Pokedex",
			"156 new Pokemon introduced in Generation V",
			domain.DataSourceTypeScheduled,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/unova",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 24,
			},
			nil,
		},
		{
			"should create a webhook data source",
			"tenant-1",
			"Paldea Pokedex",
			"Generation IX Pokemon from the Paldea region",
			domain.DataSourceTypeWebhook,
			domain.WebhookDataSourceAttributes{
				URL:     "https://example.com/paldea",
				Method:  "GET",
				Headers: map[string]string{"Authorization": "Bearer token"},
			},
			nil,
		},
		{
			"should create a data source with an empty description",
			"tenant-1",
			"Johto Pokedex",
			"",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "totodile", Label: "Totodile"},
				},
			},
			nil,
		},
		{
			"should yield an error when the name is empty",
			"tenant-1",
			"",
			"some description",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "treecko", Label: "Treecko"},
				},
			},
			errors.New("validation error"),
		},
		{
			"should yield an error when the type is invalid",
			"tenant-1",
			"Hoenn Pokedex",
			"135 new Pokemon introduced in Generation III",
			"invalid",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "mudkip", Label: "Mudkip"},
				},
			},
			domain.ErrInvalidSourceType,
		},
		{
			"should yield an error when the attributes for the type are invalid",
			"tenant-1",
			"Sinnoh Pokedex",
			"107 new Pokemon introduced in Generation IV",
			domain.DataSourceTypeStatic,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/sinnoh",
				Method:        "GET",
				IntervalHours: 12,
			},
			domain.ErrInvalidSourceTypeAttributes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := domain.NewDataSource(tt.tenantID, tt.dataSourceName, tt.description, tt.dataSourceType, tt.attributes)

			// Assert.
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if errors.Is(tt.wantErr, domain.ErrInvalidSourceType) || errors.Is(tt.wantErr, domain.ErrInvalidSourceTypeAttributes) {
					if !errors.Is(gotErr, tt.wantErr) {
						t.Errorf("expected %v but got %v", tt.wantErr, gotErr)
					}
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

			if got.ID == "" {
				t.Errorf("expected non-empty ID")
			}

			if got.TenantID != tt.tenantID {
				t.Errorf("expected tenant ID %q but got %q", tt.tenantID, got.TenantID)
			}

			if got.Name != tt.dataSourceName {
				t.Errorf("expected name %q but got %q", tt.dataSourceName, got.Name)
			}

			if got.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, got.Description)
			}

			if got.Type != tt.dataSourceType {
				t.Errorf("expected type %q but got %q", tt.dataSourceType, got.Type)
			}

			if got.CreatedAt.IsZero() {
				t.Errorf("expected non-zero CreatedAt")
			}
		})
	}
}

func TestHydrateDataSource(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		id             domain.DataSourceID
		tenantID       domain.TenantID
		dataSourceName string
		description    string
		dataSourceType domain.DataSourceType
		attributes     domain.DataSourceAttributes
		createdAt      time.Time
		updatedAt      time.Time
	}{
		{
			"should hydrate a static data source",
			"ds-1",
			"tenant-1",
			"Kanto Pokedex",
			"The original 151 Pokemon from the Kanto region",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "pikachu", Label: "Pikachu"},
					{Value: "charizard", Label: "Charizard"},
				},
			},
			now,
			now.Add(24 * time.Hour),
		},
		{
			"should hydrate a scheduled data source",
			"ds-2",
			"tenant-1",
			"Unova Pokedex",
			"156 new Pokemon introduced in Generation V",
			domain.DataSourceTypeScheduled,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/unova",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 24,
			},
			now,
			now.Add(48 * time.Hour),
		},
		{
			"should hydrate a webhook data source with empty optional fields",
			"ds-3",
			"tenant-1",
			"Paldea Pokedex",
			"",
			domain.DataSourceTypeWebhook,
			domain.WebhookDataSourceAttributes{
				URL:     "https://example.com/paldea",
				Method:  "GET",
				Headers: map[string]string{"Authorization": "Bearer token"},
			},
			now,
			time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := domain.HydrateDataSource(tt.id, tt.tenantID, tt.dataSourceName, tt.description, tt.dataSourceType, tt.attributes, tt.createdAt, tt.updatedAt)

			// Assert.
			if got == nil {
				t.Errorf("expected data source but got nil")
				return
			}

			if got.ID != tt.id {
				t.Errorf("expected ID %q but got %q", tt.id, got.ID)
			}

			if got.TenantID != tt.tenantID {
				t.Errorf("expected tenant ID %q but got %q", tt.tenantID, got.TenantID)
			}

			if got.Name != tt.dataSourceName {
				t.Errorf("expected name %q but got %q", tt.dataSourceName, got.Name)
			}

			if got.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, got.Description)
			}

			if got.Type != tt.dataSourceType {
				t.Errorf("expected type %q but got %q", tt.dataSourceType, got.Type)
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

func TestDataSource_UpdateAttributes(t *testing.T) {
	now := time.Now()

	newScheduledDS := func() *domain.DataSource {
		return domain.HydrateDataSource(
			"ds-1", "tenant-1", "Unova Pokedex", "156 new Pokemon introduced in Generation V",
			domain.DataSourceTypeScheduled,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/unova",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 24,
			},
			now, time.Time{},
		)
	}

	tests := []struct {
		name    string
		ds      *domain.DataSource
		attr    domain.DataSourceAttributes
		wantErr error
	}{
		{
			"should update attributes for a scheduled data source",
			newScheduledDS(),
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/unova",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 12,
				Attempts:      1,
			},
			nil,
		},
		{
			"should yield an error when the attributes do not match the type",
			newScheduledDS(),
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "pikachu", Label: "Pikachu"},
				},
			},
			domain.ErrInvalidSourceTypeAttributes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			gotErr := tt.ds.UpdateAttributes(tt.attr)

			// Assert.
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("expected %v but got %v", tt.wantErr, gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if tt.ds.Attributes == nil {
				t.Errorf("expected attributes to be updated")
			}

			if tt.ds.UpdatedAt.IsZero() {
				t.Errorf("expected non-zero UpdatedAt")
			}
		})
	}
}

func TestDataSource_Update(t *testing.T) {
	now := time.Now()

	newStaticDS := func() *domain.DataSource {
		return domain.HydrateDataSource(
			"ds-1", "tenant-1", "Kanto Pokedex", "The original 151 Pokemon",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{{Value: "pikachu", Label: "Pikachu"}},
			},
			now, time.Time{},
		)
	}

	tests := []struct {
		name           string
		ds             *domain.DataSource
		dataSourceName string
		description    string
		dataSourceType domain.DataSourceType
		attributes     domain.DataSourceAttributes
		wantErr        error
	}{
		{
			"should update a data source",
			newStaticDS(),
			"Johto Pokedex",
			"100 new Pokemon introduced in Generation II",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "totodile", Label: "Totodile"},
					{Value: "cyndaquil", Label: "Cyndaquil"},
				},
			},
			nil,
		},
		{
			"should update a data source with an empty description",
			newStaticDS(),
			"Hoenn Pokedex",
			"",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "treecko", Label: "Treecko"},
				},
			},
			nil,
		},
		{
			"should update a data source type",
			newStaticDS(),
			"Sinnoh Pokedex",
			"107 new Pokemon introduced in Generation IV",
			domain.DataSourceTypeScheduled,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/sinnoh",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 12,
			},
			nil,
		},
		{
			"should yield an error when the name is empty",
			newStaticDS(),
			"",
			"some description",
			domain.DataSourceTypeStatic,
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "piplup", Label: "Piplup"},
				},
			},
			errors.New("validation error"),
		},
		{
			"should yield an error when the type is invalid",
			newStaticDS(),
			"Alola Pokedex",
			"Generation VII Pokemon from the Alola region",
			"invalid",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "rowlet", Label: "Rowlet"},
				},
			},
			domain.ErrInvalidSourceType,
		},
		{
			"should yield an error when the attributes do not match the type",
			newStaticDS(),
			"Galar Pokedex",
			"Generation VIII Pokemon from the Galar region",
			domain.DataSourceTypeStatic,
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/galar",
				Method:        "GET",
				IntervalHours: 24,
			},
			domain.ErrInvalidSourceTypeAttributes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			originalName := tt.ds.Name

			// Act.
			gotErr := tt.ds.Update(tt.dataSourceName, tt.description, tt.dataSourceType, tt.attributes)

			// Assert.
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
					return
				}

				if errors.Is(tt.wantErr, domain.ErrInvalidSourceType) || errors.Is(tt.wantErr, domain.ErrInvalidSourceTypeAttributes) {
					if !errors.Is(gotErr, tt.wantErr) {
						t.Errorf("expected %v but got %v", tt.wantErr, gotErr)
					}
				}

				if tt.ds.Name != originalName {
					t.Errorf("expected name to remain %q but got %q", originalName, tt.ds.Name)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			if tt.ds.Name != tt.dataSourceName {
				t.Errorf("expected name %q but got %q", tt.dataSourceName, tt.ds.Name)
			}

			if tt.ds.Description != tt.description {
				t.Errorf("expected description %q but got %q", tt.description, tt.ds.Description)
			}

			if tt.ds.Type != tt.dataSourceType {
				t.Errorf("expected type %q but got %q", tt.dataSourceType, tt.ds.Type)
			}

			if tt.ds.UpdatedAt.IsZero() {
				t.Errorf("expected non-zero UpdatedAt")
			}
		})
	}
}
