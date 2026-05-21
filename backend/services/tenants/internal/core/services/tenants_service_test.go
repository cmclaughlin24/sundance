package services

import (
	"context"
	"errors"
	"testing"

	"sundance/backend/pkg/common"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

func TestTenantsService_Find(t *testing.T) {
	tests := []struct {
		name    string
		want    []*domain.Tenant
		wantErr error
	}{
		{
			"should yield a list of tenants",
			[]*domain.Tenant{
				{Name: "Star Fox 64", Description: "Featured the iconic phrase 'Do a barrel roll' and introduced the Rumble Pak accessory"},
				{Name: "Star Fox Adventures", Description: "Originally developed as Dinosaur Planet by Rare before becoming a Star Fox title"},
			},
			nil,
		},
		{
			"should yield an empty list of tenants",
			[]*domain.Tenant{},
			nil,
		},
		{
			"should yield an error if when the repository returns an error",
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := tenantsService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					findFn: func(_ context.Context) ([]*domain.Tenant, error) {
						return tt.want, tt.wantErr
					},
				},
			}

			// Act.
			got, gotErr := s.Find(context.Background())

			// Assert.
			if tt.wantErr != gotErr {
				t.Errorf("expected error %v but got %v", tt.wantErr, gotErr)
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("expected %d tenants but got %d", len(tt.want), len(got))
				return
			}

			for idx, want := range tt.want {
				if want.Name != got[idx].Name || want.Description != got[idx].Description {
					t.Errorf("expected %v but got %v", want, got[idx])
					break
				}
			}
		})
	}
}

func TestTenantsService_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		id      domain.TenantID
		want    *domain.Tenant
		wantErr error
	}{
		{
			"should yield a tenant",
			"star-fox-command",
			&domain.Tenant{Name: "Star Fox Command", Description: "Released in 2006 for the Nintendo DS with stylus-controlled flight and multiple branching storylines"},
			nil,
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			"star-fox-guard",
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the repository returns an error",
			"star-fox-64-3d",
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := tenantsService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					findByIdFn: func(_ context.Context, _ domain.TenantID) (*domain.Tenant, error) {
						return tt.want, tt.wantErr
					},
				},
			}

			// Act.
			got, gotErr := s.FindByID(context.Background(), tt.id)

			// Assert.
			if tt.wantErr != gotErr {
				t.Errorf("expected error %v but got %v", tt.wantErr, gotErr)
				return
			}

			if tt.want == nil {
				return
			}

			if tt.want.Name != got.Name || tt.want.Description != got.Description {
				t.Errorf("expected %v but got %v", tt.want, got)
				return
			}
		})
	}
}

func TestTenantsService_Create(t *testing.T) {
	tests := []struct {
		name    string
		command *ports.CreateTenantCommand
		want    *domain.Tenant
		wantErr error
	}{
		{
			"should create a tenant",
			ports.NewCreateTenantCommand("Star Fox 64", "Released in 1997 for the Nintendo 64, featuring Fox McCloud and his team"),
			&domain.Tenant{Name: "Star Fox 64", Description: "Released in 1997 for the Nintendo 64, featuring Fox McCloud and his team"},
			nil,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewCreateTenantCommand("", ""),
			nil,
			errors.New("validation error"),
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewCreateTenantCommand("Star Fox Assault", "Released in 2005 for the GameCube, featured on-foot combat missions"),
			nil,
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := tenantsService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					upsertFn: func(_ context.Context, _ *domain.Tenant) (*domain.Tenant, error) {
						return tt.want, tt.wantErr
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

			if tt.want == nil {
				return
			}

			if tt.want.Name != got.Name || tt.want.Description != got.Description {
				t.Errorf("expected %v but got %v", tt.want, got)
				return
			}
		})
	}
}

func TestTenantsService_Update(t *testing.T) {
	tests := []struct {
		name    string
		command *ports.UpdateTenantCommand
		want    *domain.Tenant
		wantErr error
	}{
		{
			"should update a tenant",
			ports.NewUpdateTenantCommand(domain.TenantID("star-fox-1"), "Star Fox", "Originally released in 1993 for the SNES, it was the first game to use the Super FX chip for 3D polygon graphics"),
			&domain.Tenant{Name: "Star Fox", Description: "Originally released in 1993 for the SNES, it was the first game to use the Super FX chip for 3D polygon graphics"},
			nil,
		},
		{
			"should yield an error when the command is invalid",
			ports.NewUpdateTenantCommand(domain.TenantID(""), "", ""),
			nil,
			errors.New("validation error"),
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			ports.NewUpdateTenantCommand(domain.TenantID("star-fox-2"), "Star Fox 2", "Completed in 1995 but shelved until its official release on the SNES Classic Edition in 2017"),
			nil,
			common.ErrNotFound,
		},
		{
			"should yield an error when the repository returns an error",
			ports.NewUpdateTenantCommand(domain.TenantID("star-fox-zero"), "Star Fox Zero", "Released in 2016 for the Wii U with motion controls and a dual-screen gameplay mechanic"),
			&domain.Tenant{Name: "Star Fox Zero", Description: "Released in 2016 for the Wii U with motion controls and a dual-screen gameplay mechanic"},
			errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := tenantsService{
				logger: logger,
				tenantsRepository: &mockTenantsRepository{
					findByIdFn: func(_ context.Context, _ domain.TenantID) (*domain.Tenant, error) {
						if errors.Is(tt.wantErr, common.ErrNotFound) {
							return nil, tt.wantErr
						}
						return tt.want, nil
					},
					upsertFn: func(_ context.Context, _ *domain.Tenant) (*domain.Tenant, error) {
						return tt.want, tt.wantErr
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

			if tt.want == nil {
				return
			}

			if tt.want.Name != got.Name || tt.want.Description != got.Description {
				t.Errorf("expected %v but got %v", tt.want, got)
				return
			}
		})
	}
}

func TestTenantsService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		id           domain.TenantID
		beginTxErr   error
		tenantExists bool
		existsErr    error
		deleteErr    error
		deleteAllErr error
		commitErr    error
		wantErr      bool
	}{
		{
			"should delete a tenant",
			"star-fox-64",
			nil,
			true,
			nil,
			nil,
			nil,
			nil,
			false,
		},
		{
			"should yield an error when the transaction fails to begin",
			"star-fox-adventures",
			errors.New("begin tx error"),
			false,
			nil,
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the existence check fails",
			"star-fox-assault",
			nil,
			false,
			errors.New("exists error"),
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an ErrNotFound when the tenant does not exist",
			"star-fox-command",
			nil,
			false,
			nil,
			nil,
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the tenant delete fails",
			"star-fox-zero",
			nil,
			true,
			nil,
			errors.New("delete error"),
			nil,
			nil,
			true,
		},
		{
			"should yield an error when the data sources delete fails",
			"star-fox-2",
			nil,
			true,
			nil,
			nil,
			errors.New("delete all error"),
			nil,
			true,
		},
		{
			"should yield an error when the commit fails",
			"star-fox-guard",
			nil,
			true,
			nil,
			nil,
			nil,
			errors.New("commit error"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := tenantsService{
				logger: logger,
				database: &mockDatabase{
					beginTxFn: func(ctx context.Context) (context.Context, error) {
						if tt.beginTxErr != nil {
							return ctx, tt.beginTxErr
						}
						return ctx, nil
					},
					rollbackTxfn: func(_ context.Context) error {
						return nil
					},
					commitTxFn: func(_ context.Context) error {
						return tt.commitErr
					},
				},
				tenantsRepository: &mockTenantsRepository{
					existsFn: func(_ context.Context, _ domain.TenantID) (bool, error) {
						return tt.tenantExists, tt.existsErr
					},
					deleteFn: func(_ context.Context, _ domain.TenantID) error {
						return tt.deleteErr
					},
				},
				dataSourcesRepository: &mockDataSourcesRepository{
					deleteAllFn: func(_ context.Context, _ domain.TenantID) error {
						return tt.deleteAllErr
					},
				},
			}

			// Act.
			gotErr := s.Delete(context.Background(), tt.id)

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
