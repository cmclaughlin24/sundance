package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

func init() {
	validate.RegisterValidation("keypath", validate.NewRegexValidator(`^[a-zA-Z_][a-zA-Z0-9_-]*(\[\*\])?(\.[a-zA-Z_][a-zA-Z0-9_-]*(\[\*\])?)*$`))
}

type CreateTagCommand struct {
	TenantID      string             `validate:"required"`
	KeyPath       string             `validate:"required,nowhitespace,keypath"`
	DisplayName   string             `validate:"required"`
	NodeType      domain.TagNodeType `validate:"required"`
	PrimitiveType *domain.TagPrimitiveType
	IsCollection  bool
}

func NewCreateTagCommand(tenantID, keyPath, displayName string, nodeType domain.TagNodeType, primitiveType *domain.TagPrimitiveType, isCollection bool) CreateTagCommand {
	return CreateTagCommand{
		TenantID:      tenantID,
		KeyPath:       keyPath,
		DisplayName:   displayName,
		NodeType:      nodeType,
		PrimitiveType: primitiveType,
		IsCollection:  isCollection,
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
