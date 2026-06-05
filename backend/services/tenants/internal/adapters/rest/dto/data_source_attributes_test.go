package dto

import (
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"testing"
	"time"
)

func Test_dataSourceAttributesToResponse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		attr     domain.DataSourceAttributes
		wantType string
	}{
		{
			"should yield a static attributes response",
			domain.StaticDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "monaco", Label: "Circuit de Monaco"},
					{Value: "silverstone", Label: "Silverstone Circuit"},
				},
			},
			"static",
		},
		{
			"should yield a scheduled attributes response",
			domain.ScheduledDataSourceAttributes{
				Data: []*domain.Lookup{
					{Value: "verstappen", Label: "Max Verstappen"},
				},
				DataSourceHTTPRequest: domain.DataSourceHTTPRequest{
					URL:        "https://example.com/standings",
					Method:     "GET",
					Headers:    map[string]string{"Authorization": "Bearer fia-token"},
					ValueField: "driverId",
					LabelField: "driverName",
				},
				IntervalHours:  24,
				ExpirationDate: now.Add(48 * time.Hour),
			},
			"scheduled",
		},
		{
			"should yield a webhook attributes response",
			domain.WebhookDataSourceAttributes{
				DataSourceHTTPRequest: domain.DataSourceHTTPRequest{
					URL:        "https://example.com/timing",
					Method:     "POST",
					Headers:    map[string]string{"Authorization": "Bearer fia-token"},
					ValueField: "driverId",
					LabelField: "driverName",
				},
				RequiredKeys: []string{"driver"},
			},
			"webhook",
		},
		{
			"should yield a data lake attributes response",
			domain.DataLakeDataSourceAttributes{
				Query:        "SELECT value, label FROM laps WHERE driver = @driver",
				RequiredKeys: []string{"driver"},
				OptionalKeys: []string{"season"},
				Catalog:      "f1",
				Schema:       "telemetry",
				ValueField:   "value",
				LabelField:   "label",
				TimeoutMs:    7000,
			},
			"data-lake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := dataSourceAttributesToResponse(tt.attr)

			// Assert.
			if got == nil {
				t.Errorf("expected response but got nil")
				return
			}

			switch tt.wantType {
			case "static":
				resp, ok := got.(staticDataSourceAttributesResponse)
				if !ok {
					t.Errorf("expected staticDataSourceAttributesResponse but got %T", got)
					return
				}

				input := tt.attr.(domain.StaticDataSourceAttributes)

				if len(resp.Data) != len(input.Data) {
					t.Errorf("expected %d lookups but got %d", len(input.Data), len(resp.Data))
				}
			case "scheduled":
				resp, ok := got.(scheduledDataSourceAttributesResponse)
				if !ok {
					t.Errorf("expected scheduledDataSourceAttributesResponse but got %T", got)
					return
				}

				input := tt.attr.(domain.ScheduledDataSourceAttributes)

				if resp.URL != input.URL {
					t.Errorf("expected URL %q but got %q", input.URL, resp.URL)
				}

				if resp.Method != input.Method {
					t.Errorf("expected Method %q but got %q", input.Method, resp.Method)
				}

				if resp.ValueField != input.ValueField {
					t.Errorf("expected ValueField %q but got %q", input.ValueField, resp.ValueField)
				}

				if resp.LabelField != input.LabelField {
					t.Errorf("expected LabelField %q but got %q", input.LabelField, resp.LabelField)
				}

				if resp.IntervalHours != input.IntervalHours {
					t.Errorf("expected IntervalHours %v but got %v", input.IntervalHours, resp.IntervalHours)
				}

				if len(resp.Data) != len(input.Data) {
					t.Errorf("expected %d lookups but got %d", len(input.Data), len(resp.Data))
				}
			case "webhook":
				resp, ok := got.(webhookDataSourceAttributesResponse)
				if !ok {
					t.Errorf("expected webhookDataSourceAttributesResponse but got %T", got)
					return
				}

				input := tt.attr.(domain.WebhookDataSourceAttributes)

				if resp.URL != input.URL {
					t.Errorf("expected URL %q but got %q", input.URL, resp.URL)
				}

				if resp.Method != input.Method {
					t.Errorf("expected Method %q but got %q", input.Method, resp.Method)
				}

				if len(resp.RequiredKeys) != len(input.RequiredKeys) {
					t.Errorf("expected %d required keys but got %d", len(input.RequiredKeys), len(resp.RequiredKeys))
				}

				if resp.ValueField != input.ValueField {
					t.Errorf("expected ValueField %q but got %q", input.ValueField, resp.ValueField)
				}

				if resp.LabelField != input.LabelField {
					t.Errorf("expected LabelField %q but got %q", input.LabelField, resp.LabelField)
				}
			case "data-lake":
				resp, ok := got.(dataLakeDataSourceAttributesResponse)
				if !ok {
					t.Errorf("expected dataLakeDataSourceAttributesResponse but got %T", got)
					return
				}

				input := tt.attr.(domain.DataLakeDataSourceAttributes)

				if resp.Query != input.Query {
					t.Errorf("expected Query %q but got %q", input.Query, resp.Query)
				}

				if resp.Catalog != input.Catalog {
					t.Errorf("expected Catalog %q but got %q", input.Catalog, resp.Catalog)
				}

				if resp.Schema != input.Schema {
					t.Errorf("expected Schema %q but got %q", input.Schema, resp.Schema)
				}

				if resp.ValueField != input.ValueField {
					t.Errorf("expected ValueField %q but got %q", input.ValueField, resp.ValueField)
				}

				if resp.LabelField != input.LabelField {
					t.Errorf("expected LabelField %q but got %q", input.LabelField, resp.LabelField)
				}

				if resp.TimeoutMs != input.TimeoutMs {
					t.Errorf("expected TimeoutMs %d but got %d", input.TimeoutMs, resp.TimeoutMs)
				}
			}
		})
	}
}

func TestRequestToDataSourceAttributes(t *testing.T) {
	tests := []struct {
		name           string
		dataSourceType domain.DataSourceType
		raw            any
		wantErr        bool
	}{
		{
			"should parse static attributes",
			domain.DataSourceTypeStatic,
			map[string]any{
				"data": []map[string]string{
					{"value": "monaco", "label": "Circuit de Monaco"},
					{"value": "spa", "label": "Circuit de Spa-Francorchamps"},
				},
			},
			false,
		},
		{
			"should parse scheduled attributes",
			domain.DataSourceTypeScheduled,
			map[string]any{
				"url":           "https://example.com/standings",
				"method":        "GET",
				"headers":       map[string]string{"Authorization": "Bearer fia-token"},
				"intervalHours": 24,
			},
			false,
		},
		{
			"should parse webhook attributes",
			domain.DataSourceTypeWebhook,
			map[string]any{
				"url":     "https://example.com/timing",
				"method":  "POST",
				"headers": map[string]string{"Authorization": "Bearer fia-token"},
			},
			false,
		},
		{
			"should parse data lake attributes",
			domain.DataSourceTypeDataLake,
			map[string]any{
				"query":        "SELECT value, label FROM laps WHERE driver = @driver",
				"requiredKeys": []string{"driver"},
				"optionalKeys": []string{"season"},
				"catalog":      "f1",
				"schema":       "telemetry",
				"valueField":   "value",
				"labelField":   "label",
				"limit":        250,
				"timeoutMs":    7000,
			},
			false,
		},
		{
			"should yield an error when the type is empty",
			"",
			map[string]any{},
			true,
		},
		{
			"should yield an error when the type is unknown",
			"unknown",
			map[string]any{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got, gotErr := RequestToDataSourceAttributes(tt.dataSourceType, tt.raw)

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
				t.Errorf("expected attributes but got nil")
				return
			}

			switch tt.dataSourceType {
			case domain.DataSourceTypeStatic:
				if _, ok := got.(domain.StaticDataSourceAttributes); !ok {
					t.Errorf("expected StaticDataSourceAttributes but got %T", got)
				}
			case domain.DataSourceTypeScheduled:
				if _, ok := got.(domain.ScheduledDataSourceAttributes); !ok {
					t.Errorf("expected ScheduledDataSourceAttributes but got %T", got)
				}
			case domain.DataSourceTypeWebhook:
				if _, ok := got.(domain.WebhookDataSourceAttributes); !ok {
					t.Errorf("expected WebhookDataSourceAttributes but got %T", got)
				}
			case domain.DataSourceTypeDataLake:
				if _, ok := got.(domain.DataLakeDataSourceAttributes); !ok {
					t.Errorf("expected DataLakeDataSourceAttributes but got %T", got)
				}
			}
		})
	}
}

func TestRequestToDataSourceAttributes_ErrorWrapping(t *testing.T) {
	// Act.
	_, gotErr := RequestToDataSourceAttributes("unknown", map[string]any{})

	// Assert.
	if gotErr == nil {
		t.Errorf("expected error but got nil")
		return
	}

	if !errors.Is(gotErr, ErrDataSourceAttrParse) {
		t.Errorf("expected error to wrap ErrDataSourceAttrParse but got %v", gotErr)
	}
}

