package services

import (
	"context"
	"errors"
	"log/slog"

	"sundance/backend/pkg/common"
	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type formsService struct {
	logger             *slog.Logger
	database           database.Database
	formsRepository    ports.FormsRepository
	versionsRepository ports.FormVersionRepository
}

func NewFormsService(logger *slog.Logger, repository *ports.Repository) ports.FormsAPI {
	return &formsService{
		logger:             logger,
		database:           repository.Database,
		formsRepository:    repository.Forms,
		versionsRepository: repository.FormVersions,
	}
}

func (s *formsService) Find(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
	s.logger.DebugContext(ctx, "listing forms", "tenant_id", query.TenantID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	forms, err := s.formsRepository.Find(ctx, &ports.FormFilters{
		TenantID: query.TenantID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve forms", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return forms, nil
}

func (s *formsService) FindByID(ctx context.Context, query *ports.FindFormByIDQuery) (*domain.Form, error) {
	s.logger.DebugContext(ctx, "finding form", "tenant_id", query.TenantID, "form_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form find failed; invalid query", "tenant_id", query.TenantID, "form_id", query.ID, "error", err)
		return nil, err
	}

	form, err := s.formsRepository.FindByID(ctx, query.ID)
	if err != nil {
		s.logFindFormByIDError(ctx, err, query.ID)
		return nil, err
	}

	if form.TenantID != query.TenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", query.TenantID, "form_id", query.ID)
		return nil, common.ErrUnauthorized
	}

	return form, nil
}

func (s *formsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	s.logger.DebugContext(ctx, "creating form", "tenant_id", command.TenantID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form creation failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	form, err := domain.NewForm(command.TenantID, command.Name, command.Description)
	if err != nil {
		s.logger.WarnContext(ctx, "form creation failed; domain invariant violation", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist form", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "form created", "tenant_id", command.TenantID, "form_id", form.ID)

	return form, nil
}

func (s *formsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	s.logger.DebugContext(ctx, "updating form", "tenant_id", command.TenantID, "form_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form update failed; invalid command", "tenant_id", command.TenantID, "form_id", command.ID, "error", err)
		return nil, err
	}

	form, err := s.formsRepository.FindByID(ctx, command.ID)
	if err != nil {
		s.logFindFormByIDError(ctx, err, command.ID)
		return nil, err
	}

	if form.TenantID != command.TenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", command.TenantID, "form_id", command.ID)
		return nil, common.ErrUnauthorized
	}

	if err := form.Update(command.Name, command.Description); err != nil {
		s.logger.WarnContext(ctx, "form update failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.ID, "error", err)
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist form", "tenant_id", command.TenantID, "form_id", command.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "form updated", "tenant_id", command.TenantID, "form_id", command.ID)

	return form, nil
}

func (s *formsService) Delete(ctx context.Context, command *ports.RemoveFormCommand) error {
	s.logger.DebugContext(ctx, "deleting form", "tenant_id", command.TenantID, "form_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form deletion failed; invalid command", "tenant_id", command.TenantID, "form_id", command.ID, "error", err)
		return err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.ID); err != nil {
		return err
	}

	hasActive, err := s.hasActiveVersion(ctx, command.ID)
	if err != nil {
		return err
	} else if hasActive {
		s.logger.WarnContext(ctx, "form deletion failed; form has active version", "tenant_id", command.TenantID, "form_id", command.ID)
		return domain.ErrFormHasActiveVersion
	}

	if err := s.formsRepository.Delete(ctx, command.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete form", "tenant_id", command.TenantID, "form_id", command.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "form deleted", "tenant_id", command.TenantID, "form_id", command.ID)

	return nil
}

func (s *formsService) FindVersions(ctx context.Context, query *ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "listing versions", "tenant_id", query.TenantID, "form_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version listing failed; invalid query", "tenant_id", query.TenantID, "form_id", query.ID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.ID); err != nil {
		return nil, err
	}

	versions, err := s.versionsRepository.Find(ctx, query.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve versions", "tenant_id", query.TenantID, "form_id", query.ID, "error", err)
		return nil, err
	}

	return versions, nil
}

func (s *formsService) FindVersion(ctx context.Context, query *ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "finding version", "tenant_id", query.TenantID, "form_id", query.ID, "version_id", query.VersionID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version find failed; invalid query", "tenant_id", query.TenantID, "form_id", query.ID, "version_id", query.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.ID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, query.ID, query.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, query.ID, query.VersionID)
		return nil, err
	}

	return version, nil
}

func (s *formsService) CreateVersion(ctx context.Context, command *ports.CreateFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "creating version", "tenant_id", command.TenantID, "form_id", command.FormID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version creation failed; invalid command", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	versionNum, err := s.versionsRepository.FindNextVersionNumber(txCtx, command.FormID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to find next version number", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	version, err := domain.NewFormVersion(command.FormID, versionNum, domain.FormVersionStatusDraft)
	if err != nil {
		s.logger.WarnContext(ctx, "version creation failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	if err := version.AddPages(command.Pages...); err != nil {
		s.logger.WarnContext(ctx, "version creation failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(txCtx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit version creation transaction", "tenant_id", command.TenantID, "form_id", command.FormID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version created", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", version.ID)

	return version, nil
}

func (s *formsService) UpdateVersion(ctx context.Context, command *ports.UpdateFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "updating version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version update failed; invalid command", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, command.FormID, command.VersionID)
		return nil, err
	}

	if err := version.ReplacePages(command.Pages...); err != nil {
		s.logger.WarnContext(ctx, "version update failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version updated", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	return version, nil
}

func (s *formsService) PublishVersion(ctx context.Context, command *ports.PublishFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "publishing version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version publish failed; invalid command", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, command.FormID, command.VersionID)
		return nil, err
	}

	if err := version.Publish(command.UserID); err != nil {
		s.logger.WarnContext(ctx, "version publish failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version published", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	return version, nil
}

func (s *formsService) RetireVersion(ctx context.Context, command *ports.RetireFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "retiring version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version retire failed; invalid command", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, command.FormID, command.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, command.FormID, command.VersionID)
		return nil, err
	}

	if err := version.Retire(command.UserID); err != nil {
		s.logger.WarnContext(ctx, "version retire failed; domain invariant violation", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version retired", "tenant_id", command.TenantID, "form_id", command.FormID, "version_id", command.VersionID)

	return version, nil
}

func (s *formsService) isValidAccess(ctx context.Context, tenantID string, formID domain.FormID) error {
	form, err := s.formsRepository.FindByID(ctx, formID)

	if err != nil {
		s.logFindFormByIDError(ctx, err, formID)
		return err
	}

	if form.TenantID != tenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", tenantID, "form_id", formID)
		return common.ErrUnauthorized
	}

	return nil
}

func (s *formsService) hasActiveVersion(ctx context.Context, id domain.FormID) (bool, error) {
	versions, err := s.versionsRepository.Find(ctx, id)

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve versions", "form_id", id, "error", err)
		return true, err
	}

	for _, v := range versions {
		if v.Status == domain.FormVersionStatusActive {
			return true, nil
		}
	}

	return false, nil
}

func (s *formsService) logFindFormByIDError(ctx context.Context, err error, formID domain.FormID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "form not found", "form_id", formID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve form", "form_id", formID, "error", err)
	}
}

func (s *formsService) logFindVersionByIDError(ctx context.Context, err error, formID domain.FormID, versionID domain.FormVersionID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "version not found", "form_id", formID, "version_id", versionID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve version", "form_id", formID, "version_id", versionID, "error", err)
	}
}
