package services

import (
	"context"
	"errors"
	"log/slog"
	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type canonicalTagService struct {
	logger     *slog.Logger
	repository ports.CanonicalTagRepository
}

func NewCanonicalTagService(logger *slog.Logger, repository *ports.Repository) ports.CanonicalTagAPI {
	return &canonicalTagService{
		logger:     logger,
		repository: repository.CanonicalTags,
	}
}

func (s *canonicalTagService) Find(ctx context.Context, query ports.FindCanonicalTagsQuery) ([]*domain.CanonicalTag, error) {
	s.logger.DebugContext(ctx, "listing canonical tags", "tenant_id", "")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tags, err := s.repository.Find(ctx, ports.CanonicalTagFilters{
		TenantID: query.TenantID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve canonical tags", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return tags, nil
}

func (s *canonicalTagService) FindById(ctx context.Context, query ports.FindByIDQuery[domain.CanonicalTagID]) (*domain.CanonicalTag, error) {
	s.logger.DebugContext(ctx, "finding canonical tag", "tenant_id", query.TenantID, "canonical_tag_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "find canonical tag failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	tag, err := s.repository.FindByID(ctx, query.ID)
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

func (s *canonicalTagService) Create(ctx context.Context, command ports.CreateCanonicalTagCommand) (*domain.CanonicalTag, error) {
	s.logger.DebugContext(ctx, "creating canoonical tag", "tenant_id", command.TenantID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "canonical tag creation failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	tag, err := domain.NewCanonicalTag(command.TenantID, command.Key, command.DisplayName)
	if err != nil {
		s.logger.WarnContext(ctx, "canonical tag creation failed; domain invariant violation", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	tag, err = s.repository.Upsert(ctx, tag)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist canonical tag", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "canonical tag created", "tenant_id", command.TenantID, "canonical_tag_id", tag.ID)

	return tag, nil
}

func (s *canonicalTagService) Delete(ctx context.Context, command any) error {
	return nil
}

func (s *canonicalTagService) logFindByIDError(ctx context.Context, err error, tagID domain.CanonicalTagID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "canonical tag not found", "canonical_tag_id", tagID)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve canonical tag", "canonical_tag_id", tagID, "error", err)
	}
}
