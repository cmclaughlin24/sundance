package services

import (
	"context"
	"log/slog"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type tagMappingValidator struct {
	logger             *slog.Logger
	tagsRepository     ports.TagsRepository
	versionsRepository ports.TagVersionsRepository
}

func newTagMappingValidator(logger *slog.Logger, repository *ports.Repository) *tagMappingValidator {
	return &tagMappingValidator{
		logger:             logger,
		tagsRepository:     repository.Tags,
		versionsRepository: repository.TagVersions,
	}
}

func (v *tagMappingValidator) validate(ctx context.Context, pages []*domain.Page) error {
	return nil
}
