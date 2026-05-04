package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type TenantsService struct {
	logger                *log.Logger
	database              database.Database
	tenantsRepository     ports.TenantsRepository
	dataSourcesRepository ports.DataSourcesRepository
}

func NewTenantsService(logger *log.Logger, repository *ports.Repository) ports.TenantsService {
	return &TenantsService{
		logger:                logger,
		database:              repository.Database,
		tenantsRepository:     repository.Tenants,
		dataSourcesRepository: repository.DataSources,
	}
}

func (s *TenantsService) Find(ctx context.Context) ([]*domain.Tenant, error) {
	return s.tenantsRepository.Find(ctx)
}

func (s *TenantsService) FindByID(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return s.tenantsRepository.FindByID(ctx, id)
}

func (s *TenantsService) Create(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	t, err := domain.NewTenant(command.Name, command.Description)
	if err != nil {
		return nil, err
	}

	return s.tenantsRepository.Upsert(ctx, t)
}

func (s *TenantsService) Update(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	tenant, err := s.tenantsRepository.FindByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	if err := tenant.Update(command.Name, command.Description); err != nil {
		return nil, err
	}

	return s.tenantsRepository.Upsert(ctx, tenant)
}

func (s *TenantsService) Delete(ctx context.Context, id domain.TenantID) error {
	txCtx, err := s.database.BeginTx(ctx)

	if err != nil {
		return err
	}

	defer s.database.RollbackTx(txCtx)
	exists, err := s.tenantsRepository.Exists(txCtx, id)

	if err != nil {
		return err
	}

	if !exists {
		return common.ErrNotFound
	}

	if err := s.tenantsRepository.Delete(txCtx, id); err != nil {
		return err
	}

	if err := s.dataSourcesRepository.DeleteAll(txCtx, id); err != nil {
		return err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		return err
	}

	return nil
}
