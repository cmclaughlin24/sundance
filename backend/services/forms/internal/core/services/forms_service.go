package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type FormsService struct {
	logger             *log.Logger
	database           database.Database
	formsRepository    ports.FormsRepository
	versionsRepository ports.VersionRepository
}

func NewFormsService(logger *log.Logger, repository *ports.Repository) ports.FormsService {
	return &FormsService{
		logger:          logger,
		database:        repository.Database,
		formsRepository: repository.Forms,
	}
}

func (s *FormsService) Find(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	return s.formsRepository.Find(ctx, &ports.FormFilters{
		TenantID: query.TenantID,
	})
}

func (s *FormsService) FindByID(ctx context.Context, query *ports.FindFormsByIDQuery) (*domain.Form, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	form, err := s.formsRepository.FindByID(ctx, query.FormID)

	if err != nil {
		return nil, err
	}

	if form.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return form, nil
}

func (s *FormsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	form, err := domain.NewForm(command.TenantID, command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	form, err := s.formsRepository.FindByID(ctx, command.ID)

	if err != nil {
		return nil, err
	}

	if form.TenantID != command.TenantID {
		return nil, common.ErrUnauthorized
	}

	if err := form.Update(command.Name, command.Description); err != nil {
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) Delete(ctx context.Context, command *ports.RemoveFormCommand) error {
	if err := validate.ValidateStruct(command); err != nil {
		return err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.ID); err != nil {
		return err
	}

	hasActive, err := s.hasActiveVersion(ctx, command.ID)

	if err != nil {
		return err
	}

	if hasActive {
		return domain.ErrFormHasActiveVersion
	}

	return s.formsRepository.Delete(ctx, command.ID)
}

func (s *FormsService) FindVersions(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.FormID); err != nil {
		return nil, err
	}

	return s.versionsRepository.Find(ctx, query.FormID)
}

func (s *FormsService) FindVersion(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.FormID); err != nil {
		return nil, err
	}

	return s.versionsRepository.FindByID(ctx, query.FormID, query.VersionID)
}

func (s *FormsService) CreateVersion(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)

	if err != nil {
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	versionNum, err := s.versionsRepository.FindNextVersionNumber(txCtx, command.FormID)

	if err != nil {
		return nil, err
	}

	version, err := domain.NewVersion(command.FormID, versionNum, domain.VersionStatusDraft)

	if err != nil {
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(txCtx, version)

	if err != nil {
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) UpdateVersion(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.UpdatePages(command.Pages...); err != nil {
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) PublishVersion(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Publish(command.UserID); err != nil {
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) RetireVersion(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Retire(command.UserID); err != nil {
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) isValidAccess(ctx context.Context, tenantID string, formID domain.FormID) error {
	form, err := s.formsRepository.FindByID(ctx, formID)

	if err != nil {
		return err
	}

	if form.TenantID != tenantID {
		return common.ErrUnauthorized
	}

	return nil
}

func (s *FormsService) hasActiveVersion(ctx context.Context, id domain.FormID) (bool, error) {
	versions, err := s.versionsRepository.Find(ctx, id)

	if err != nil {
		return true, err
	}

	for _, v := range versions {
		if v.Status == domain.VersionStatusActive {
			return true, err
		}
	}

	return false, nil
}
