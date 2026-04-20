package services

import (
	"context"
	"log"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type FormsService struct {
	logger          *log.Logger
	database        database.Database
	formsRepository ports.FormsRepository
	baseService
}

func NewFormsService(logger *log.Logger, repository *ports.Repository) *FormsService {
	return &FormsService{
		logger:          logger,
		database:        repository.Database,
		formsRepository: repository.Forms,
	}
}

func (s *FormsService) Find(ctx context.Context) ([]*domain.Form, error) {
	return s.formsRepository.Find(ctx)
}

func (s *FormsService) FindById(ctx context.Context, query *ports.FindByIDQuery) (*domain.Form, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	form, err := s.formsRepository.FindById(ctx, query.FormID)

	if err != nil {
		return nil, err
	}

	if form.TenantID != tenantID {
		return nil, common.ErrUnauthorized
	}

	return form, nil
}

func (s *FormsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	form, err := domain.NewForm(domain.FormID(""), tenantID, command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	form, err = s.formsRepository.Create(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, command.ID); err != nil {
		return nil, err
	}

	form, err := s.formsRepository.FindById(ctx, command.ID)

	if err != nil {
		return nil, err
	}

	if err := form.Update(command.Name, command.Description); err != nil {
		return nil, err
	}

	form, err = s.formsRepository.Update(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) FindVersions(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, query.FormID); err != nil {
		return nil, err
	}

	return s.formsRepository.FindVersions(ctx, query.FormID)
}

func (s *FormsService) FindVersion(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, query.FormID); err != nil {
		return nil, err
	}

	return s.formsRepository.FindVersion(ctx, query.FormID, query.VersionID)
}

func (s *FormsService) CreateVersion(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, command.FormID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)

	if err != nil {
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	versionNum, err := s.formsRepository.FindNextVersionNumber(txCtx, command.FormID)

	if err != nil {
		return nil, err
	}

	version, err := domain.NewVersion("", command.FormID, versionNum, domain.VersionStatusDraft)

	if err != nil {
		return nil, err
	}

	version, err = s.formsRepository.CreateVersion(txCtx, version)

	if err != nil {
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) UpdateVersion(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.formsRepository.FindVersion(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.UpdatePages(command.Pages...); err != nil {
		return nil, err
	}

	version, err = s.formsRepository.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) PublishVersion(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.formsRepository.FindVersion(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Publish(command.UserID, time.Now()); err != nil {
		return nil, err
	}

	version, err = s.formsRepository.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) RetireVersion(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
	tenantID, err := s.getTenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, tenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.formsRepository.FindVersion(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Retire(command.UserID, time.Now()); err != nil {
		return nil, err
	}

	version, err = s.formsRepository.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) isValidAccess(ctx context.Context, tenantID string, formID domain.FormID) error {
	form, err := s.formsRepository.FindById(ctx, formID)

	if err != nil {
		return err
	}

	if form.TenantID != tenantID {
		return common.ErrUnauthorized
	}

	return nil
}
