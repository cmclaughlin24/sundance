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

type tagsService struct {
	logger             *slog.Logger
	database           database.Database
	tagsRepository     ports.TagsRepository
	versionsRepository ports.TagVersionsRepository
}

func NewTagsService(logger *slog.Logger, repository *ports.Repository) ports.TagsAPI {
	return &tagsService{
		logger:             logger,
		database:           repository.Database,
		tagsRepository:     repository.Tags,
		versionsRepository: repository.TagVersions,
	}
}

func (s *tagsService) Find(ctx context.Context, query ports.FindTagsQuery) ([]*domain.Tag, error) {
	s.logger.DebugContext(ctx, "listing tags", "tenant_id", query.TenantID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tags, err := s.tagsRepository.Find(ctx, ports.TagFilters{
		TenantID: query.TenantID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve tags", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return tags, nil
}

func (s *tagsService) FindById(ctx context.Context, query ports.FindByIDQuery[domain.TagID]) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "finding tag", "tenant_id", query.TenantID, "tag_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "find tag failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tag, err := s.tagsRepository.FindByID(ctx, query.ID)
	if err != nil {
		s.logFindTagByIDError(ctx, err, query.ID)
		return nil, err
	}

	if tag.TenantID != query.TenantID {
		s.logger.WarnContext(ctx, "unauthorized tag access", "tenant_id", query.TenantID, "tag_id", query.ID)
		return nil, common.ErrUnauthorized
	}

	return tag, nil
}

func (s *tagsService) Create(ctx context.Context, cmd commands.CreateTagCommand) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "creating tag", "tenant_id", cmd.TenantID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag creation failed; invalid command", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	tag, err := domain.NewTag(cmd.TenantID, cmd.KeyPath, cmd.DisplayName, cmd.NodeType, cmd.PrimitiveType)
	if err != nil {
		s.logger.WarnContext(ctx, "tag creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	tag, err = s.tagsRepository.Upsert(ctx, tag)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist tag", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tag created", "tenant_id", cmd.TenantID, "tag_id", tag.ID)

	return tag, nil
}

func (s *tagsService) Update(ctx context.Context, cmd commands.UpdateTagCommand) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "updating tag", "tenant_id", cmd.TenantID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag update failed; invalid command", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	tag, err := s.tagsRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		s.logFindTagByIDError(ctx, err, cmd.ID)
		return nil, err
	}

	if tag.TenantID != cmd.TenantID {
		s.logger.WarnContext(ctx, "unauthorized tag access", "tenant_id", cmd.TenantID, "tag_id", cmd.ID)
		return nil, common.ErrUnauthorized
	}

	if err := tag.Update(cmd.DisplayName); err != nil {
		s.logger.WarnContext(ctx, "tag update failed; domain invariant violation", "tenant_id", cmd.TenantID, "tag_id", cmd.ID, "error", err)
		return nil, err
	}

	tag, err = s.tagsRepository.Upsert(ctx, tag)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist tag", "tenant_id", cmd.TenantID, "tag_id", cmd.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tag updated", "tenant_id", cmd.TenantID, "tag_id", cmd.ID)

	return tag, nil
}

func (s *tagsService) Delete(ctx context.Context, cmd commands.DeleteCommand[domain.TagID]) error {
	s.logger.DebugContext(ctx, "deleting tag", "tenant_id", cmd.TenantID, "tag_id", cmd.ID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag deletion failed; invalid command", "tenant_id", cmd.TenantID, "tag_id", cmd.ID, "error", err)
		return err
	}

	versions, err := s.versionsRepository.Find(ctx, ports.TagVersionFilters{TagID: cmd.ID})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve versions", "tenant_id", cmd.TenantID, "tag_id", cmd.ID, "error", err)
		return err
	} else if len(versions) > 0 {
		s.logger.WarnContext(ctx, "tag deletion failed; tag has active version", "tenant_id", cmd.TenantID, "tag_id", cmd.ID)
		return domain.ErrTagHasVersions
	}

	if err := s.tagsRepository.Delete(ctx, cmd.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete tag", "tenant_id", cmd.TenantID, "tag_id", cmd.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "tag deleted", "tenant_id", cmd.TenantID, "tag_id", cmd.ID)

	return nil
}

func (s *tagsService) FindVersions(ctx context.Context, query ports.FindTagVersionsQuery) ([]*domain.TagVersion, error) {
	s.logger.DebugContext(ctx, "listing versions", "tenant_id", query.TenantID, "canonical_version_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version listing failed; invalid query", "tenant_id", query.TenantID, "canonical_version_id", query.ID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.ID); err != nil {
		return nil, err
	}

	versions, err := s.versionsRepository.Find(ctx, ports.TagVersionFilters{
		TagID: query.ID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve versions", "tenant_id", query.TenantID, "canonical_version_id", query.ID, "error", err)
		return nil, err
	}

	return versions, nil
}

func (s *tagsService) FindVersion(ctx context.Context, query ports.FindTagVersionQuery) (*domain.TagVersion, error) {
	s.logger.DebugContext(ctx, "finding version", "tenant_id", query.TenantID, "tag_id", query.ID, "version_id", query.VersionID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version find failed; invalid query", "tenant_id", query.TenantID, "tag_id", query.ID, "version_id", query.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, query.TenantID, query.ID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, query.VersionID)
	if err != nil {
		s.logFindTagVersionByIDError(ctx, err, query.ID, query.VersionID)
		return nil, err
	}

	return version, nil
}

func (s *tagsService) CreateVersion(ctx context.Context, cmd commands.CreateTagVersionCommand) (*domain.TagVersion, error) {
	s.logger.DebugContext(ctx, "creating version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "version creation failed; invalid command", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.TagID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	versionNum, err := s.versionsRepository.FindNextVersionNumber(txCtx, cmd.TagID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to find next version number", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	version, err := domain.NewTagVersion(cmd.TagID, versionNum)
	if err != nil {
		s.logger.WarnContext(ctx, "version creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(txCtx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit version creation transaction", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "version created", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", version.ID)

	return version, nil
}

func (s *tagsService) PublishVersion(ctx context.Context, cmd commands.TransitionTagVersionCommand) (*domain.TagVersion, error) {
	s.logger.DebugContext(ctx, "transitioning tag version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag version transition failed; invalid command", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.TagID); err != nil {
		return nil, err
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	defer s.database.RollbackTx(txCtx)

	if err := s.deprecateActiveVersions(txCtx, cmd.TagID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(txCtx, cmd.VersionID)
	if err != nil {
		s.logFindTagVersionByIDError(ctx, err, cmd.TagID, cmd.VersionID)
		return nil, err
	}

	if err := version.Publish(); err != nil {
		s.logger.WarnContext(ctx, "tag version transition failed; domain invariant violation", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(txCtx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist tag version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit publish transaction", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tag version transitioned", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID)

	return version, nil
}

func (s *tagsService) DeprecateVersion(ctx context.Context, cmd commands.TransitionTagVersionCommand) (*domain.TagVersion, error) {
	return s.transitionVersion(ctx, cmd, func(tv *domain.TagVersion) error {
		return tv.Deprecate()
	})
}

func (s *tagsService) RetireVersion(ctx context.Context, cmd commands.TransitionTagVersionCommand) (*domain.TagVersion, error) {
	return s.transitionVersion(ctx, cmd, func(tv *domain.TagVersion) error {
		return tv.Retire()
	})
}

func (s *tagsService) transitionVersion(ctx context.Context, cmd commands.TransitionTagVersionCommand, transition func(*domain.TagVersion) error) (*domain.TagVersion, error) {
	s.logger.DebugContext(ctx, "transitioning tag version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "tag version transition failed; invalid command", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	if err := s.isValidAccess(ctx, cmd.TenantID, cmd.TagID); err != nil {
		return nil, err
	}

	version, err := s.versionsRepository.FindByID(ctx, cmd.VersionID)
	if err != nil {
		s.logFindTagVersionByIDError(ctx, err, cmd.TagID, cmd.VersionID)
		return nil, err
	}

	if err := transition(version); err != nil {
		s.logger.WarnContext(ctx, "tag version transition failed; domain invariant violation", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	version, err = s.versionsRepository.Upsert(ctx, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist tag version", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "tag version transitioned", "tenant_id", cmd.TenantID, "tag_id", cmd.TagID, "version_id", cmd.VersionID)

	return version, nil
}

func (s *tagsService) deprecateActiveVersions(ctx context.Context, tagID domain.TagID) error {
	versions, err := s.versionsRepository.Find(ctx, ports.TagVersionFilters{
		TagID:    tagID,
		Statuses: []domain.TagStatus{domain.TagStatusActive},
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve tag versions", "tag_id", tagID, "error", err)
		return err
	}

	for _, v := range versions {
		if err := v.Deprecate(); err != nil {
			s.logger.WarnContext(ctx, "tag version deprecation failed; domain invariant violation", "tag_id", v.TagID, "version_id", v.ID, "error", err)
			return err
		}

		v, err = s.versionsRepository.Upsert(ctx, v)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to persist tag version", "tag_id", tagID, "version_id", v.ID, "error", err)
			return err
		}
	}

	return nil
}

func (s *tagsService) logFindTagByIDError(ctx context.Context, err error, tagID domain.TagID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "tag not found", "tag_id", tagID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve tag", "tag_id", tagID, "error", err)
	}
}

func (s *tagsService) logFindTagVersionByIDError(ctx context.Context, err error, tagID domain.TagID, versionID domain.TagVersionID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "tag not found", "tag_id", tagID, "tag_version_id", versionID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve tag", "tag_id", tagID, "tag_version_id", versionID, "error", err)
	}
}

func (s *tagsService) isValidAccess(ctx context.Context, tenantID string, id domain.TagID) error {
	form, err := s.tagsRepository.FindByID(ctx, id)

	if err != nil {
		s.logFindTagByIDError(ctx, err, id)
		return err
	}

	if form.TenantID != tenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", tenantID, "tag_id", id)
		return common.ErrUnauthorized
	}

	return nil
}
