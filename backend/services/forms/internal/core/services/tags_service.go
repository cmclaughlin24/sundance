package services

import (
	"context"
	"errors"
	"log/slog"
	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type tagsService struct {
	logger             *slog.Logger
	tagsRepository     ports.TagsRepository
	versionsRepository ports.TagVersionsRepository
}

func NewTagsService(logger *slog.Logger, repository *ports.Repository) ports.TagsAPI {
	return &tagsService{
		logger:             logger,
		tagsRepository:     repository.Tags,
		versionsRepository: repository.TagVersions,
	}
}

func (s *tagsService) Find(ctx context.Context, query ports.FindTagsQuery) ([]*domain.Tag, error) {
	s.logger.DebugContext(ctx, "listing canonical tags", "tenant_id", "")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tags, err := s.tagsRepository.Find(ctx, ports.TagFilters{
		TenantID: query.TenantID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve canonical tags", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return tags, nil
}

func (s *tagsService) FindById(ctx context.Context, query ports.FindByIDQuery[domain.TagID]) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "finding canonical tag", "tenant_id", query.TenantID, "canonical_tag_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "find canonical tag failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tag, err := s.tagsRepository.FindByID(ctx, query.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, query.ID)
		return nil, err
	}

	if tag.TenantID != query.TenantID {
		s.logger.WarnContext(ctx, "unauthorized canonical tag access", "tenant_id", query.TenantID, "canonical_tag_id", query.ID)
		return nil, common.ErrUnauthorized
	}

	return tag, nil
}

func (s *tagsService) Create(ctx context.Context, command ports.CreateTagCommand) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "creating canonical tag", "tenant_id", command.TenantID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag creation failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	tag, err := domain.NewTag(command.TenantID, command.Key, command.DisplayName)
	if err != nil {
		s.logger.WarnContext(ctx, "canonical tag creation failed; domain invariant violation", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	tag, err = s.tagsRepository.Upsert(ctx, tag)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist canonical tag", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "canonical tag created", "tenant_id", command.TenantID, "canonical_tag_id", tag.ID)

	return tag, nil
}

func (s *tagsService) Update(ctx context.Context, command ports.UpdateTagCommand) (*domain.Tag, error) {
	s.logger.DebugContext(ctx, "updating canonical tag", "tenant_id", command.TenantID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag update failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	tag, err := s.tagsRepository.FindByID(ctx, command.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, command.ID)
		return nil, err
	}

	if tag.TenantID != command.TenantID {
		s.logger.WarnContext(ctx, "unauthorized canonical tag access", "tenant_id", command.TenantID, "canonical_tag_id", command.ID)
		return nil, common.ErrUnauthorized
	}

	if err := tag.Update(command.DisplayName); err != nil {
		s.logger.WarnContext(ctx, "canonical tag update failed; domain invariant violation", "tenant_id", command.TenantID, "canonical_tag_id", command.ID, "error", err)
		return nil, err
	}

	tag, err = s.tagsRepository.Upsert(ctx, tag)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist canonical tag", "tenant_id", command.TenantID, "canonical_tag_id", command.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "canonical tag updated", "tenant_id", command.TenantID, "canonical_tag_id", command.ID)

	return tag, nil
}

func (s *tagsService) Delete(ctx context.Context, command ports.DeleteCommand[domain.TagID]) error {
	s.logger.DebugContext(ctx, "deleting canonical tag", "tenant_id", command.TenantID, "canonical_tag_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag deletion failed; invalid command", "tenant_id", command.TenantID, "canonical_tag_id", command.ID, "error", err)
		return err
	}

	// FIXME: A canonical tag should not be deletable if it has ever had an active version to ensure audit history can be maintained.

	if err := s.tagsRepository.Delete(ctx, command.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete canonical tag", "tenant_id", command.TenantID, "canonical_tag_id", command.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "canonical tag deleted", "tenant_id", command.TenantID, "canonical_tag_id", command.ID)

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
	return nil, nil
}

func (s *tagsService) logFindByIDError(ctx context.Context, err error, tagID domain.TagID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "canonical tag not found", "canonical_tag_id", tagID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve canonical tag", "canonical_tag_id", tagID, "error", err)
	}
}

func (s *tagsService) isValidAccess(ctx context.Context, tenantID string, id domain.TagID) error {
	form, err := s.tagsRepository.FindByID(ctx, id)

	if err != nil {
		s.logFindByIDError(ctx, err, id)
		return err
	}

	if form.TenantID != tenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", tenantID, "canonical_tag_id", id)
		return common.ErrUnauthorized
	}

	return nil
}
