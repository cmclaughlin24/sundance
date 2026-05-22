package domain_test

import (
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func TestScheduledDataSourceAttributes_RefreshData(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		data          []*domain.Lookup
		intervalHours float64
	}{
		{
			"should refresh data with new lookups",
			[]*domain.Lookup{
				{Value: "pikachu", Label: "Pikachu"},
				{Value: "eevee", Label: "Eevee"},
			},
			24,
		},
		{
			"should refresh data with an empty list",
			[]*domain.Lookup{},
			12,
		},
		{
			"should replace existing data",
			[]*domain.Lookup{
				{Value: "bulbasaur", Label: "Bulbasaur"},
			},
			48,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			attr := domain.ScheduledDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "charmander", Label: "Charmander"},
					{Value: "squirtle", Label: "Squirtle"},
				},
				IntervalHours: tt.intervalHours,
				Attempts:      3,
			}

			// Act.
			attr.RefreshData(tt.data)

			// Assert.
			if len(attr.Data) != len(tt.data) {
				t.Errorf("expected %d lookups but got %d", len(tt.data), len(attr.Data))
				return
			}

			for idx, want := range tt.data {
				if attr.Data[idx].Value != want.Value || attr.Data[idx].Label != want.Label {
					t.Errorf("expected %v but got %v", want, attr.Data[idx])
					break
				}
			}

			if attr.ExpirationDate.IsZero() {
				t.Errorf("expected non-zero ExpirationDate")
			}

			if !attr.ExpirationDate.After(now) {
				t.Errorf("expected ExpirationDate to be after %v but got %v", now, attr.ExpirationDate)
			}

			if attr.Attempts != 0 {
				t.Errorf("expected Attempts to be reset to 0 but got %d", attr.Attempts)
			}
		})
	}
}

func TestScheduledDataSourceAttributes_RecordAttempt(t *testing.T) {
	tests := []struct {
		name            string
		initialAttempts int
		wantAttempts    int
	}{
		{
			"should increment attempts from zero",
			0,
			1,
		},
		{
			"should increment attempts from existing value",
			2,
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			attr := domain.ScheduledDataSourceAttributes{
				Attempts: tt.initialAttempts,
			}

			// Act.
			attr.RecordAttempt()

			// Assert.
			if attr.Attempts != tt.wantAttempts {
				t.Errorf("expected Attempts %d but got %d", tt.wantAttempts, attr.Attempts)
			}
		})
	}
}

func TestGetDataSourceAttributes_Static(t *testing.T) {
	tests := []struct {
		name    string
		attr    domain.DataSourceAttributes
		wantErr bool
	}{
		{
			"should return static attributes",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "pikachu", Label: "Pikachu"},
				},
			},
			false,
		},
		{
			"should yield an error for scheduled attributes",
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/pokemon",
				Method:        "GET",
				IntervalHours: 24,
			},
			true,
		},
		{
			"should yield an error for nil attributes",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := domain.GetDataSourceAttributes[domain.StaticDataSourceAttributes](tt.attr)

			// Assert.
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
				}
				if !errors.Is(gotErr, domain.ErrDataSourceAttributeMismatch) {
					t.Errorf("expected ErrDataSourceAttributeMismatch but got %v", gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			input, _ := tt.attr.(domain.StaticDataSourceAttributes)
			if len(got.Data) != len(input.Data) {
				t.Errorf("expected %d lookups but got %d", len(input.Data), len(got.Data))
			}
		})
	}
}

func TestGetDataSourceAttributes_Scheduled(t *testing.T) {
	tests := []struct {
		name    string
		attr    domain.DataSourceAttributes
		wantErr bool
	}{
		{
			"should return scheduled attributes",
			domain.ScheduledDataSourceAttributes{
				URL:           "https://example.com/pokemon",
				Method:        "GET",
				Headers:       map[string]string{"Authorization": "Bearer token"},
				IntervalHours: 24,
			},
			false,
		},
		{
			"should yield an error for static attributes",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "charmander", Label: "Charmander"},
				},
			},
			true,
		},
		{
			"should yield an error for nil attributes",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := domain.GetDataSourceAttributes[domain.ScheduledDataSourceAttributes](tt.attr)

			// Assert.
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
				}
				if !errors.Is(gotErr, domain.ErrDataSourceAttributeMismatch) {
					t.Errorf("expected ErrDataSourceAttributeMismatch but got %v", gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			input, _ := tt.attr.(domain.ScheduledDataSourceAttributes)
			if got.URL != input.URL {
				t.Errorf("expected URL %q but got %q", input.URL, got.URL)
			}
		})
	}
}

func TestGetDataSourceAttributes_Webhook(t *testing.T) {
	tests := []struct {
		name    string
		attr    domain.DataSourceAttributes
		wantErr bool
	}{
		{
			"should return webhook attributes",
			domain.WebhookDataSourceAttributes{
				URL:     "https://example.com/pokemon",
				Method:  "GET",
				Headers: map[string]string{"Authorization": "Bearer token"},
			},
			false,
		},
		{
			"should yield an error for static attributes",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "squirtle", Label: "Squirtle"},
				},
			},
			true,
		},
		{
			"should yield an error for nil attributes",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := domain.GetDataSourceAttributes[domain.WebhookDataSourceAttributes](tt.attr)

			// Assert.
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("expected error but got nil")
				}
				if !errors.Is(gotErr, domain.ErrDataSourceAttributeMismatch) {
					t.Errorf("expected ErrDataSourceAttributeMismatch but got %v", gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("expected no error but got %v", gotErr)
				return
			}

			input, _ := tt.attr.(domain.WebhookDataSourceAttributes)
			if got.URL != input.URL {
				t.Errorf("expected URL %q but got %q", input.URL, got.URL)
			}
		})
	}
}

