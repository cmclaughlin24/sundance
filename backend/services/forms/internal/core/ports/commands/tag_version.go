package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type CreateTagVersionCommand struct {
	TenantID string         `validate:"required"`
	TagID    domain.TagID   `validate:"required"`
	Type     domain.TagType `validate:"required"`
}

func NewCreateTagVersionCommand(tenantID string, tagID domain.TagID, tagType domain.TagType) CreateTagVersionCommand {
	return CreateTagVersionCommand{
		TenantID: tenantID,
		TagID:    tagID,
		Type:     tagType,
	}
}

func (c CreateTagVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type TransitionTagVersionCommand struct {
	TenantID  string              `validate:"required"`
	TagID     domain.TagID        `validate:"required"`
	VersionID domain.TagVersionID `validate:"required"`
}

func NewTransitionTagVersionCommand(tenantID string, tagID domain.TagID, versionID domain.TagVersionID) TransitionTagVersionCommand {
	return TransitionTagVersionCommand{
		TenantID:  tenantID,
		TagID:     tagID,
		VersionID: versionID,
	}
}

func (c TransitionTagVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}
