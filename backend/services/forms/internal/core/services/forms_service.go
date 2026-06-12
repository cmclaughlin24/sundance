package services

import (
	"context"
	"errors"
	"log/slog"

	"sundance/backend/pkg/common"
	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/ports/commands"
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

func (s *formsService) Find(ctx context.Context, query ports.FindFormsQuery) ([]*domain.Form, error) {
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

func (s *formsService) FindByID(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
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

func (s *formsService) Create(ctx context.Context, cmd commands.CreateFormCommand) (*domain.Form, error) {
	s.logger.DebugContext(ctx, "creating form", "tenant_id", cmd.TenantID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form creation failed; invalid command", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	form, err := domain.NewForm(cmd.TenantID, cmd.Name, cmd.Description)
	if err != nil {
		s.logger.WarnContext(ctx, "form creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist form", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "form created", "tenant_id", cmd.TenantID, "form_id", form.ID)

	return form, nil
}

func (s *formsService) Update(ctx context.Context, cmd commands.UpdateFormCommand) (*domain.Form, error) {
	s.logger.DebugContext(ctx, "updating form", "tenant_id", cmd.TenantID, "form_id", cmd.ID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form update failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.ID, "error", err)
		return nil, err
	}

	form, err := s.formsRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		s.logFindFormByIDError(ctx, err, cmd.ID)
		return nil, err
	}

	if form.TenantID != cmd.TenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", cmd.TenantID, "form_id", cmd.ID)
		return nil, common.ErrUnauthorized
	}

	if err := form.Update(cmd.Name, cmd.Description); err != nil {
		s.logger.WarnContext(ctx, "form update failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.ID, "error", err)
		return nil, err
	}

	form, err = s.formsRepository.Upsert(ctx, form)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist form", "tenant_id", cmd.TenantID, "form_id", cmd.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "form updated", "tenant_id", cmd.TenantID, "form_id", cmd.ID)

	return form, nil
}

func (s *formsService) Delete(ctx context.Context, cmd commands.DeleteCommand[domain.FormID]) error {
	s.logger.DebugContext(ctx, "deleting form", "tenant_id", cmd.TenantID, "form_id", cmd.ID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "form deletion failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.ID, "error", err)
		return err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.ID); err != nil {
		return err
	}

	hasActive, err := s.hasActiveVersion(ctx, cmd.ID)
	if err != nil {
		return err
	} else if hasActive {
		s.logger.WarnContext(ctx, "form deletion failed; form has active version", "tenant_id", cmd.TenantID, "form_id", cmd.ID)
		return domain.ErrFormHasActiveVersion
	}

	if err := s.formsRepository.Delete(ctx, cmd.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete form", "tenant_id", cmd.TenantID, "form_id", cmd.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "form deleted", "tenant_id", cmd.TenantID, "form_id", cmd.ID)

	return nil
}

func (s *formsService) FindVersions(ctx context.Context, query ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
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

func (s *formsService) FindVersion(ctx context.Context, query ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "finding version", "tenant_id", query.TenantID, "form_id", query.ID, "version_id", query.VersionID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version find failed; invalid query", "tenant_id", query.TenantID, "form_id", query.ID, "version_id", query.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.ID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, query.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, query.ID, query.VersionID)
		return nil, err
	}

	return version, nil
}

func (s *formsService) CreateVersion(ctx context.Context, cmd *commands.CreateFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "creating version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version creation failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.FormID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	versionNum, err := s.versionsRepository.FindNextVersionNumber(txCtx, cmd.FormID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to find next version number", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	version, err := domain.NewFormVersion(cmd.FormID, versionNum, domain.FormVersionStatusDraft)
	if err != nil {
		s.logger.WarnContext(ctx, "version creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	if err := version.AddPages(cmd.Pages...); err != nil {
		s.logger.WarnContext(ctx, "version creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(txCtx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit version creation transaction", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version created", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", version.ID)

	return version, nil
}

func (s *formsService) UpdateVersion(ctx context.Context, cmd *commands.UpdateFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "updating version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version update failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, cmd.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, cmd.FormID, cmd.VersionID)
		return nil, err
	}

	if err := version.ReplacePages(cmd.Pages...); err != nil {
		s.logger.WarnContext(ctx, "version update failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version updated", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

	return version, nil
}

func (s *formsService) PublishVersion(ctx context.Context, cmd commands.PublishFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "publishing version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version publish failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, cmd.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, cmd.FormID, cmd.VersionID)
		return nil, err
	}

	if err := version.Publish(cmd.UserID); err != nil {
		s.logger.WarnContext(ctx, "version publish failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version published", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

	return version, nil
}

func (s *formsService) RetireVersion(ctx context.Context, cmd commands.RetireFormVersionCommand) (*domain.FormVersion, error) {
	s.logger.DebugContext(ctx, "retiring version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version retire failed; invalid command", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.FormID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, cmd.VersionID)
	if err != nil {
		s.logFindVersionByIDError(ctx, err, cmd.FormID, cmd.VersionID)
		return nil, err
	}

	if err := version.Retire(cmd.UserID); err != nil {
		s.logger.WarnContext(ctx, "version retire failed; domain invariant violation", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version retired", "tenant_id", cmd.TenantID, "form_id", cmd.FormID, "version_id", cmd.VersionID)

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
