package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
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
	return nil, nil
}

func (s *TenantsService) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return nil, nil
}

func (s *TenantsService) Create(ctx context.Context, command ports.CreateTenantCommand) (*domain.Tenant, error) {
	return nil, nil
}

func (s *TenantsService) Update(ctx context.Context, command ports.UpdateTenantCommand) (*domain.Tenant, error) {
	return nil, nil
}

func (s *TenantsService) Remove(ctx context.Context, id domain.TenantID) error {
	return nil
}
