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
			ports.NewFindDataSourceJobsQuery(0),
			[]*domain.DataSource{
				{TenantID: "tenant-1", Name: "Source 1"},
				{TenantID: "tenant-1", Name: "Source 2"},
			},
			nil,
		},
		{
			"should yield an empty list of data sources",
			ports.NewFindDataSourceJobsQuery(0),
			[]*domain.DataSource{},
			nil,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewFindDataSourceJobsQuery(0),
			nil,
			errors.New("repository error"),
		},
		{
			"should yield an error when the query is invalid",
			ports.NewFindDataSourceJobsQuery(-1),
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
			URL:           "https://example.com/pokemon-types",
			Method:        "GET",
			Headers:       map[string]string{"Authorization": "Bearer token"},
			IntervalHours: 24,
		},
	}

	lookups := []*domain.Lookup{
		{Value: "fire", Label: "Fire"},
		{Value: "water", Label: "Water"},
	}

	tests := []struct {
		name           string
		command        *ports.ProcessDataSourceJobCommand
		fetchLookupsFn func(context.Context, string, string, map[string]string) ([]*domain.Lookup, error)
		upsertErr      error
		wantErr        bool
	}{
		{
			"should process a data source job",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return lookups, nil
			},
			nil,
			false,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewProcessDataSourceJobCommand(nil),
			nil,
			nil,
			true,
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
		},
		{
			"should yield an error when fetching lookups fails",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return nil, errors.New("fetch error")
			},
			nil,
			true,
		},
		{
			"should yield an error when the repository fails to persist",
			ports.NewProcessDataSourceJobCommand(scheduledDS),
			func(_ context.Context, _, _ string, _ map[string]string) ([]*domain.Lookup, error) {
				return lookups, nil
			},
			errors.New("repository error"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesJobService{
				logger: logger,
				repository: &mockDataSourcesRepository{
					upsertFn: func(_ context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
						if tt.upsertErr != nil {
							return nil, tt.upsertErr
						}
						return ds, nil
					},
				},
				client: &mockLookupClient{
					fetchLookupsFn: func(ctx context.Context, method, url string, headers map[string]string) ([]*domain.Lookup, error) {
						if tt.fetchLookupsFn != nil {
							return tt.fetchLookupsFn(ctx, method, url, headers)
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
		})
	}
}
