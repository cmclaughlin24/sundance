package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type TenantsService struct {
	logger            *log.Logger
	tenantsRepository ports.TenantsRepository
}

func NewTenantsService(logger *log.Logger, repository *ports.Repository) ports.TenantsService {
	return &TenantsService{
		logger:            logger,
		tenantsRepository: repository.Tenants,
	}
}

func (s *TenantsService) Find(ctx context.Context) ([]*domain.Tenant, error) {
	return s.tenantsRepository.Find(ctx)
}

func (s *TenantsService) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return s.tenantsRepository.FindById(ctx, id)
}

func (s *TenantsService) Create(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	t, err := domain.NewTenant(command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantsRepository.Upsert(ctx, t)

	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantsService) Update(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	tenant, err := s.tenantsRepository.FindById(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	if err := tenant.Update(command.Name, command.Description); err != nil {
		return nil, err
	}

	tenant, err = s.tenantsRepository.Upsert(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *TenantsService) Remove(ctx context.Context, id domain.TenantID) error {
	exists, err := s.tenantsRepository.Exists(ctx, id)

	if err != nil {
		return err
	}

	if !exists {
		return common.ErrNotFound
	}

	return s.tenantsRepository.Remove(ctx, id)
}
