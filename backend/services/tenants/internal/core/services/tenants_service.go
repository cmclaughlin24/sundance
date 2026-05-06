package services

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type TenantsService struct {
	logger                *slog.Logger
	database              database.Database
	tenantsRepository     ports.TenantsRepository
	dataSourcesRepository ports.DataSourcesRepository
}

func NewTenantsService(logger *slog.Logger, repository *ports.Repository) ports.TenantsService {
	return &TenantsService{
		logger:                logger,
		database:              repository.Database,
		tenantsRepository:     repository.Tenants,
		dataSourcesRepository: repository.DataSources,
	}
}

func (s *TenantsService) Find(ctx context.Context) ([]*domain.Tenant, error) {
	s.logger.DebugContext(ctx, "listing tenants")

	tenants, err := s.tenantsRepository.Find(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve tenants", "error", err)
		return nil, err
	}

	return tenants, nil
}

func (s *TenantsService) FindByID(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	s.logger.DebugContext(ctx, "finding tenant", "tenant_id", id)

	t, err := s.tenantsRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TenantsService) Create(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
	s.logger.DebugContext(ctx, "creating tenant")

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tenant creation failed; invalid command", "error", err)
		return nil, err
	}

	t, err := domain.NewTenant(command.Name, command.Description)
	if err != nil {
		s.logger.WarnContext(ctx, "tenant creation failed; domain invariant violation", "error", err)
		return nil, err
	}

	t, err = s.tenantsRepository.Upsert(ctx, t)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist tenant", "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tenant created", "tenant_id", t.ID)

	return t, nil
}

func (s *TenantsService) Update(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
	s.logger.DebugContext(ctx, "updating tenant", "tenant_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tenant update failed; invalid command", "error", err)
		return nil, err
	}

	tenant, err := s.tenantsRepository.FindByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	if err := tenant.Update(command.Name, command.Description); err != nil {
		s.logger.WarnContext(ctx, "tenant update failed; domain invariant violation", "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tenant updated", "tenant_id", tenant.ID)

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
