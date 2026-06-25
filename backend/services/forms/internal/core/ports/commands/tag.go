package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type CreateTagCommand struct {
	TenantID     string              `validate:"required"`
	Key          string              `validate:"required,nowhitespace"`
	DisplayName  string              `validate:"required"`
	ValueKind    domain.TagValueKind `validate:"required"`
	IsCollection bool
}

func NewCreateTagCommand(tenantID, key, displayName string, valueKind domain.TagValueKind, isCollection bool) CreateTagCommand {
	return CreateTagCommand{
		TenantID:     tenantID,
		Key:          key,
		DisplayName:  displayName,
		ValueKind:    valueKind,
		IsCollection: isCollection,
	}
}

func (c CreateTagCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateTagCommand struct {
	TenantID    string       `validate:"required"`
	ID          domain.TagID `validate:"required"`
	DisplayName string       `validate:"required"`
}

func NewUpdateTagCommand(tenantID string, id domain.TagID, displayName string) UpdateTagCommand {
	return UpdateTagCommand{
		TenantID:    tenantID,
		ID:          id,
		DisplayName: displayName,
	}
}

func (c UpdateTagCommand) Validate() error {
	return validate.ValidateStruct(c)
}
