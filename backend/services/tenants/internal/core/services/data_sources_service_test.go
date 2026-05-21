package services

import (
	"context"
	"errors"
	"testing"

	"sundance/backend/pkg/common"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

func TestDataSourcesService_Find(t *testing.T) {
	tests := []struct {
		name    string
		query   *ports.ListDataSourceQuery
		want    []*domain.DataSource
		wantErr error
	}{
		{
			"should yield a list of data sources",
			ports.NewListDataSourceQuery("tenant-1"),
			[]*domain.DataSource{
				{TenantID: "tenant-1", Name: "Source 1"},
				{TenantID: "tenant-1", Name: "Source 2"},
			},
			nil,
		},
		{
			"should yield an empty list of data sources",
			ports.NewListDataSourceQuery("tenant-1"),
			[]*domain.DataSource{},
			nil,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewListDataSourceQuery("tenant-1"),
			nil,
			errors.New("repository error"),
		},
		{
			"should yield an error when the query is invalid",
			ports.NewListDataSourceQuery(""),
			nil,
			errors.New("validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				dataSourcesRepository: &mockDataSourcesRepository{
					findFn: func(_ context.Context, _ domain.TenantID) ([]*domain.DataSource, error) {
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

func TestDataSourcesService_FindByID(t *testing.T) {
	tests := []struct {
		name         string
		query        *ports.FindDataSourceByIDQuery
		tenantExists bool
		want         *domain.DataSource
		wantErr      error
	}{
		{
			"should yield a data source",
			ports.NewFindDataSourceByIDQuery("tenant-1", "ds-1"),
			true,
			&domain.DataSource{ID: "ds-1", TenantID: "tenant-1", Name: "Source 1"},
			nil,
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			ports.NewFindDataSourceByIDQuery("tenant-1", "ds-1"),
			false,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an ErrNotFound when the data source does not exist",
			ports.NewFindDataSourceByIDQuery("tenant-1", "ds-1"),
			true,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewFindDataSourceByIDQuery("tenant-1", "ds-1"),
			true,
			nil,
			errors.New("repository error"),
		},
		{
			"should yield an error when the query is invalid",
			ports.NewFindDataSourceByIDQuery("", ""),
			true,
			nil,
			errors.New("validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					existsFn: func(_ context.Context, _ domain.TenantID) (bool, error) {
						return tt.tenantExists, nil
					},
				},
				dataSourcesRepository: &mockDataSourcesRepository{
					findByIdFn: func(_ context.Context, _ domain.TenantID, _ domain.DataSourceID) (*domain.DataSource, error) {
						return tt.want, tt.wantErr
					},
				},
			}

			// Act.
			got, gotErr := s.FindByID(context.Background(), tt.query)

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

			if tt.want == nil {
				return
			}

			if tt.want.ID != got.ID || tt.want.TenantID != got.TenantID || tt.want.Name != got.Name {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestDataSourcesService_Create(t *testing.T) {
	tests := []struct {
		name         string
		command      *ports.CreateDataSourceCommand
		tenantExists bool
		want         *domain.DataSource
		wantErr      error
	}{
		{
			"should create a data source",
			ports.NewCreateDataSourceCommand(
				"tenant-1",
				"F-Zero Tracks",
				"The original F-Zero was released in 1990 for the Super Nintendo and featured 15 tracks across three leagues",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "mute-city", Label: "Mute City"},
						{Value: "big-blue", Label: "Big Blue"},
						{Value: "port-town", Label: "Port Town"},
					},
				},
			),
			true,
			&domain.DataSource{TenantID: "tenant-1", Name: "F-Zero Tracks", Description: "The original F-Zero was released in 1990 for the Super Nintendo and featured 15 tracks across three leagues"},
			nil,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewCreateDataSourceCommand("", "", "", "", nil),
			true,
			nil,
			errors.New("validation error"),
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			ports.NewCreateDataSourceCommand(
				"tenant-1",
				"F-Zero Pilots",
				"F-Zero X on the Nintendo 64 expanded the roster to 30 playable pilots",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "captain-falcon", Label: "Captain Falcon"},
						{Value: "samurai-goroh", Label: "Samurai Goroh"},
					},
				},
			),
			false,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when there is a domain invariant violation",
			ports.NewCreateDataSourceCommand(
				"tenant-1",
				"F-Zero Machines",
				"F-Zero GX was co-developed by Sega's Amusement Vision and is considered one of the fastest racing games ever made",
				"invalid",
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "blue-falcon", Label: "Blue Falcon"},
						{Value: "fire-stingray", Label: "Fire Stingray"},
					},
				},
			),
			true,
			nil,
			domain.ErrInvalidSourceType,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewCreateDataSourceCommand(
				"tenant-1",
				"F-Zero Cups",
				"F-Zero Maximum Velocity was a launch title for the Game Boy Advance in 2001",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "knight-league", Label: "Knight League"},
						{Value: "queen-league", Label: "Queen League"},
						{Value: "king-league", Label: "King League"},
					},
				}),
			true,
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					existsFn: func(_ context.Context, _ domain.TenantID) (bool, error) {
						return tt.tenantExists, nil
					},
				},
				dataSourcesRepository: &mockDataSourcesRepository{
					upsertFn: func(_ context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
						if tt.wantErr != nil {
							return nil, tt.wantErr
						}
						return ds, nil
					},
				},
			}

			// Act.
			got, gotErr := s.Create(context.Background(), tt.command)

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

			if got == nil {
				t.Errorf("expected data source but got nil")
				return
			}

			if tt.want.TenantID != got.TenantID || tt.want.Name != got.Name || tt.want.Description != got.Description {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestDataSourcesService_Update(t *testing.T) {
	validDS := &domain.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "F-Zero Tracks",
		Type:     domain.DataSourceTypeStatic,
		Attributes: domain.StaticDataSourceAttributes{
			Data: []*domain.Lookup{
				{Value: "mute-city", Label: "Mute City"},
			},
		},
	}

	tests := []struct {
		name         string
		command      *ports.UpdateDataSourceCommand
		tenantExists bool
		findDs       *domain.DataSource
		findErr      error
		upsertErr    error
		want         *domain.DataSource
		wantErr      error
	}{
		{
			"should update a data source",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Pilots",
				"F-Zero X on the Nintendo 64 expanded the roster to 30 playable pilots",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "captain-falcon", Label: "Captain Falcon"},
						{Value: "samurai-goroh", Label: "Samurai Goroh"},
					},
				},
			),
			true,
			validDS,
			nil,
			nil,
			&domain.DataSource{TenantID: "tenant-1", Name: "F-Zero Pilots", Description: "F-Zero X on the Nintendo 64 expanded the roster to 30 playable pilots"},
			nil,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewUpdateDataSourceCommand("", "", "", "", "", nil),
			true,
			nil,
			nil,
			nil,
			nil,
			errors.New("validation error"),
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Cups",
				"F-Zero Maximum Velocity was a launch title for the Game Boy Advance in 2001",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "knight-league", Label: "Knight League"},
					},
				},
			),
			false,
			nil,
			nil,
			nil,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an ErrNotFound when the data source does not exist",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Tracks",
				"The original F-Zero was released in 1990",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "mute-city", Label: "Mute City"},
					},
				},
			),
			true,
			nil,
			common.ErrNotFound,
			nil,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the repository returns an error on find",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Tracks",
				"The original F-Zero was released in 1990",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "mute-city", Label: "Mute City"},
					},
				},
			),
			true,
			nil,
			errors.New("repository error"),
			nil,
			nil,
			errors.New("repository error"),
		},
		{
			"should yield an error when there is a domain invariant violation",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Machines",
				"F-Zero GX was co-developed by Sega's Amusement Vision",
				"invalid",
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "blue-falcon", Label: "Blue Falcon"},
					},
				},
			),
			true,
			validDS,
			nil,
			nil,
			nil,
			domain.ErrInvalidSourceType,
		},
		{
			"should yield an error when the repository fails to persist",
			ports.NewUpdateDataSourceCommand(
				"tenant-1",
				"ds-1",
				"F-Zero Tracks",
				"The original F-Zero was released in 1990",
				domain.DataSourceTypeStatic,
				domain.StaticDataSourceAttributes{
					Data: []*domain.Lookup{
						{Value: "mute-city", Label: "Mute City"},
					},
				},
			),
			true,
			validDS,
			nil,
			errors.New("repository error"),
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					existsFn: func(_ context.Context, _ domain.TenantID) (bool, error) {
						return tt.tenantExists, nil
					},
				},
				dataSourcesRepository: &mockDataSourcesRepository{
					findByIdFn: func(_ context.Context, _ domain.TenantID, _ domain.DataSourceID) (*domain.DataSource, error) {
						if tt.findErr != nil {
							return nil, tt.findErr
						}
						// Return a copy so each test case gets a fresh instance.
						cpy := *tt.findDs
						return &cpy, nil
					},
					upsertFn: func(_ context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
						if tt.upsertErr != nil {
							return nil, tt.upsertErr
						}
						return ds, nil
					},
				},
			}

			// Act.
			got, gotErr := s.Update(context.Background(), tt.command)

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

			if got == nil {
				t.Errorf("expected data source but got nil")
				return
			}

			if tt.want.TenantID != got.TenantID || tt.want.Name != got.Name || tt.want.Description != got.Description {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestDataSourcesService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		command      *ports.RemoveDataSourceCommand
		tenantExists bool
		dsExists     bool
		dsExistsErr  error
		wantErr      error
	}{
		{
			"should delete a data source",
			ports.NewRemoveDataSourceCommand("tenant-1", "ds-1"),
			true,
			true,
			nil,
			nil,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewRemoveDataSourceCommand("", ""),
			true,
			true,
			nil,
			errors.New("validation error"),
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			ports.NewRemoveDataSourceCommand("tenant-1", "ds-1"),
			false,
			true,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the existence check fails",
			ports.NewRemoveDataSourceCommand("tenant-1", "ds-1"),
			true,
			false,
			errors.New("exists error"),
			errors.New("exists error"),
		},
		{
			"should yield an ErrNotFound when the data source does not exist",
			ports.NewRemoveDataSourceCommand("tenant-1", "ds-1"),
			true,
			false,
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the repository fails to delete",
			ports.NewRemoveDataSourceCommand("tenant-1", "ds-1"),
			true,
			true,
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					existsFn: func(_ context.Context, _ domain.TenantID) (bool, error) {
						return tt.tenantExists, nil
					},
				},
				dataSourcesRepository: &mockDataSourcesRepository{
					existsFn: func(_ context.Context, _ domain.TenantID, _ domain.DataSourceID) (bool, error) {
						return tt.dsExists, tt.dsExistsErr
					},
					deleteFn: func(_ context.Context, _ domain.TenantID, _ domain.DataSourceID) error {
						return tt.wantErr
					},
				},
			}

			// Act.
			gotErr := s.Delete(context.Background(), tt.command)

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
		})
	}
}

func TestDataSourcesService_Lookup(t *testing.T) {
	lookups := []*domain.Lookup{
		{Value: "mute-city", Label: "Mute City"},
		{Value: "big-blue", Label: "Big Blue"},
	}

	staticDS := &domain.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "F-Zero Tracks",
		Type:     domain.DataSourceTypeStatic,
		Attributes: domain.StaticDataSourceAttributes{
			Data: lookups,
		},
	}

	tests := []struct {
		name      string
		query     *ports.GetDataSourceLookupsQuery
		findDs    *domain.DataSource
		findErr   error
		registry  ports.LookupStrategyRegistry
		lookupErr error
		want      []*domain.Lookup
		wantErr   bool
	}{
		{
			"should yield lookups",
			ports.NewGetDataSourceLookupsQuery("tenant-1", "ds-1"),
			staticDS,
			nil,
			ports.LookupStrategyRegistry{
				domain.DataSourceTypeStatic: &mockLookupStrategy{
					lookupFn: func(_ context.Context, _ *domain.DataSource) ([]*domain.Lookup, error) {
						return lookups, nil
					},
				},
			},
			nil,
			lookups,
			false,
		},
		{
			"should yield an error when the query is invalid",
			ports.NewGetDataSourceLookupsQuery("", ""),
			nil,
			nil,
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an ErrNotFound when the data source does not exist",
			ports.NewGetDataSourceLookupsQuery("tenant-1", "ds-1"),
			nil,
			common.ErrNotFound,
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the repository returns an error on find",
			ports.NewGetDataSourceLookupsQuery("tenant-1", "ds-1"),
			nil,
			errors.New("repository error"),
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the lookup strategy is not found",
			ports.NewGetDataSourceLookupsQuery("tenant-1", "ds-1"),
			staticDS,
			nil,
			ports.LookupStrategyRegistry{},
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the lookup execution fails",
			ports.NewGetDataSourceLookupsQuery("tenant-1", "ds-1"),
			staticDS,
			nil,
			ports.LookupStrategyRegistry{
				domain.DataSourceTypeStatic: &mockLookupStrategy{
					lookupFn: func(_ context.Context, _ *domain.DataSource) ([]*domain.Lookup, error) {
						return nil, errors.New("lookup error")
					},
				},
			},
			nil,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := dataSourcesService{
				logger: logger,
				dataSourcesRepository: &mockDataSourcesRepository{
					findByIdFn: func(_ context.Context, _ domain.TenantID, _ domain.DataSourceID) (*domain.DataSource, error) {
						return tt.findDs, tt.findErr
					},
				},
				lookupStrategies: tt.registry,
			}

			// Act.
			got, gotErr := s.Lookup(context.Background(), tt.query)

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

			if len(tt.want) != len(got) {
				t.Errorf("expected %d lookups but got %d", len(tt.want), len(got))
				return
			}

			for idx, want := range tt.want {
				if want.Value != got[idx].Value || want.Label != got[idx].Label {
					t.Errorf("expected %v but got %v", want, got[idx])
					break
				}
			}
		})
	}
}
