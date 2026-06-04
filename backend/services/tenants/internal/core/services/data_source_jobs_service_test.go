package services

import (
	"context"
	"errors"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"testing"
)

func Test_dataSourcesJobService_Find(t *testing.T) {
	tests := []struct {
		name    string
		query   *ports.FindDataSourceJobsQuery
		want    []*domain.DataSource
		wantErr error
	}{
		{
			"should yield a list of data sources",
			ports.NewFindDataSourceJobsQuery(0, 1),
			[]*domain.DataSource{
				{TenantID: "tenant-1", Name: "Source 1"},
				{TenantID: "tenant-1", Name: "Source 2"},
			},
			nil,
		},
		{
			"should yield an empty list of data sources",
			ports.NewFindDataSourceJobsQuery(0, 1),
			[]*domain.DataSource{},
			nil,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewFindDataSourceJobsQuery(0, 1),
			nil,
			errors.New("repository error"),
		},
		{
			"should yield an error when the query is invalid",
			ports.NewFindDataSourceJobsQuery(-1, 1),
			nil,
			errors.New("validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesJobService{
				logger: logger,
				repository: &mockDataSourcesRepository{
					findJobsFn: func(_ context.Context, _ *ports.FindDataSourceJobsFilter) ([]*domain.DataSource, error) {
						return tt.want, tt.wantErr
					},
				},
			}

			// Act.
			got, gotErr := s.Find(context.Background(), tt.query)

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

			if len(tt.want) != len(got) {
				t.Errorf("expected %d data sources but got %d", len(tt.want), len(got))
				return
			}

			for idx, want := range tt.want {
				if want.TenantID != got[idx].TenantID || want.Name != got[idx].Name {
					t.Errorf("expected %v but got %v", want, got[idx])
					break
				}
			}
		})
	}
}

func Test_dataSourcesJobService_Process(t *testing.T) {
	scheduledDS := &domain.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Pokemon Types",
		Type:     domain.DataSourceTypeScheduled,
		Attributes: domain.ScheduledDataSourceAttributes{
			DataSourceRequest: domain.DataSourceRequest{
				URL:        "https://example.com/pokemon-types",
				Method:     "GET",
				Headers:    map[string]string{"Authorization": "Bearer token"},
				ValueField: "value",
				LabelField: "label",
			},
			IntervalHours: 24,
		},
	}

	rows := []map[string]any{
		{"value": "fire", "label": "Fire"},
		{"value": "water", "label": "Water"},
	}

	tests := []struct {
		name             string
		command          *ports.ProcessDataSourceJobCommand
		fetchLookupsFn   func(context.Context, domain.DataSourceRequest, map[string]any) ([]map[string]any, error)
		upsertErr        error
		wantErr          bool
		wantRefreshedLen int
	}{
		{
			"should process a data source job",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _ domain.DataSourceRequest, _ map[string]any) ([]map[string]any, error) {
				return rows, nil
			},
			nil,
			false,
			2,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewProcessDataSourceJobCommand(nil),
			nil,
			nil,
			true,
			0,
		},
		{
			"should yield an error when the attribute type mismatches",
			ports.NewProcessDataSourceJobCommand(&domain.DataSource{
				ID:         "ds-2",
				TenantID:   "tenant-1",
				Name:       "Pokemon Regions",
				Type:       domain.DataSourceTypeStatic,
				Attributes: domain.StaticDataSourceAttributes{},
			}),
			nil,
			nil,
			true,
			0,
		},
		{
			"should record an attempt when fetching lookups fails",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _ domain.DataSourceRequest, _ map[string]any) ([]map[string]any, error) {
				return nil, errors.New("fetch error")
			},
			nil,
			false,
			-1,
		},
		{
			"should yield an error when the repository fails to persist",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _ domain.DataSourceRequest, _ map[string]any) ([]map[string]any, error) {
				return rows, nil
			},
			errors.New("repository error"),
			true,
			0,
		},
		{
			"should skip rows missing value or label fields",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _ domain.DataSourceRequest, _ map[string]any) ([]map[string]any, error) {
				return []map[string]any{
					{"value": "fire", "label": "Fire"},
					{"value": "no-label"},
					{"label": "No Value"},
					{"value": "weird", "label": 42},
				}, nil
			},
			nil,
			false,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			var upserted *domain.DataSource
			s := dataSourcesJobService{
				logger: logger,
				repository: &mockDataSourcesRepository{
					upsertFn: func(_ context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
						if tt.upsertErr != nil {
							return nil, tt.upsertErr
						}
						upserted = ds
						return ds, nil
					},
				},
				client: &mockLookupClient{
					fetchLookupsFn: func(ctx context.Context, request domain.DataSourceRequest, params map[string]any) ([]map[string]any, error) {
						if tt.fetchLookupsFn != nil {
							return tt.fetchLookupsFn(ctx, request, params)
						}
						return nil, nil
					},
				},
			}

			// Act.
			gotErr := s.Process(context.Background(), tt.command)

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

			if tt.wantRefreshedLen >= 0 {
				if upserted == nil {
					t.Errorf("expected upsert to be called")
					return
				}
				attr, ok := upserted.Attributes.(domain.ScheduledDataSourceAttributes)
				if !ok {
					t.Errorf("expected ScheduledDataSourceAttributes but got %T", upserted.Attributes)
					return
				}
				if len(attr.Data) != tt.wantRefreshedLen {
					t.Errorf("expected refreshed data length %d but got %d", tt.wantRefreshedLen, len(attr.Data))
				}
			}
		})
	}
}
