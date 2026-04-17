package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/validate"
)

type TenantsService struct {
	logger     *log.Logger
	repository *ports.Repository
}

func NewTenantsService(logger *log.Logger, repository *ports.Repository) *TenantsService {
	return &TenantsService{
		logger:     logger,
		repository: repository,
	}
}

func (s *TenantsService) Find(ctx context.Context) ([]*domain.Tenant, error) {
	return s.repository.Tenants.Find(ctx)
}

func (s *TenantsService) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return s.repository.Tenants.FindById(ctx, id)
}

func (s *TenantsService) Create(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	t, err := domain.NewTenant("", command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	tenant, err := s.repository.Tenants.Upsert(ctx, t)

	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantsService) Update(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	t, err := domain.NewTenant(command.ID, command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	tenant, err := s.repository.Tenants.Upsert(ctx, t)

	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantsService) Remove(ctx context.Context, id domain.TenantID) error {
	return s.repository.Tenants.Remove(ctx, id)
}
