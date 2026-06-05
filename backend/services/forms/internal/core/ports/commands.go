package ports

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type DeleteCommand[T comparable] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewDeleteCommand[T comparable](tenantID string, id T) DeleteCommand[T] {
	return DeleteCommand[T]{
		TenantID: tenantID,
		ID:       id,
	}
}

func (c DeleteCommand[T]) Validate() error {
	return validate.ValidateStruct(c)
}

type baseFormCommand struct {
	TenantID    string `validate:"required"`
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=500"`
}

type CreateFormCommand struct {
	baseFormCommand
}

func NewCreateFormCommand(tenantID, name, description string) CreateFormCommand {
	return CreateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
	}
}

func (c CreateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(tenantID string, id domain.FormID, name, description string) UpdateFormCommand {
	return UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

func (c UpdateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type baseFormVersionCommand struct {
	TenantID string        `validate:"required"`
	FormID   domain.FormID `validate:"required"`
}

type CreateFormVersionCommand struct {
	baseFormVersionCommand
	Pages []*domain.Page
}

func NewCreateFormVersionCommand(tenantID string, formID domain.FormID, pages []*domain.Page) *CreateFormVersionCommand {
	return &CreateFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		Pages: pages,
	}
}

func (c *CreateFormVersionCommand) Validate() error {
	return validate.ValidateStruct(*c)
}

type UpdateFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	Pages     []*domain.Page
}

func NewUpdateFormVersionCommand(tenantID string, id domain.FormVersionID, formID domain.FormID, pages []*domain.Page) *UpdateFormVersionCommand {
	return &UpdateFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: id,
		Pages:     pages,
	}
}

func (c *UpdateFormVersionCommand) Validate() error {
	return validate.ValidateStruct(*c)
}

type PublishFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewPublishFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) PublishFormVersionCommand {
	return PublishFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c PublishFormVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RetireFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewRetireFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) RetireFormVersionCommand {
	return RetireFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c RetireFormVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type CreateSubmissionCommand struct {
	TenantID      string                         `validate:"required"`
	FormID        domain.FormID                  `validate:"required"`
	VersionID     domain.FormVersionID           `validate:"required"`
	IdempotencyID domain.IdempotencyID           `validate:"required"`
	Values        []*domain.SubmissionFieldValue `validate:"required,min=1"`
}

func NewCreateSubmissionCommand(
	tenantID string,
	formID domain.FormID,
	versionID domain.FormVersionID,
	idempotencyID domain.IdempotencyID,
	values []*domain.SubmissionFieldValue,
) *CreateSubmissionCommand {
	return &CreateSubmissionCommand{
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		IdempotencyID: idempotencyID,
		Values:        values,
	}
}

func (c *CreateSubmissionCommand) Validate() error {
	return validate.ValidateStruct(*c)
}

type ReplaySubmissionCommand struct {
	TenantID string              `validate:"required"`
	ID       domain.SubmissionID `validate:"required"`
}

func NewReplaySubmissionCommand(tenantID string, id domain.SubmissionID) ReplaySubmissionCommand {
	return ReplaySubmissionCommand{
		TenantID: tenantID,
		ID:       id,
	}
}

func (c ReplaySubmissionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type CreateTagCommand struct {
	TenantID    string `validate:"required"`
	Key         string `validate:"required,nowhitespace"`
	DisplayName string `validate:"required"`
}

func NewCreateTagCommand(tenantID, key, displayName string) CreateTagCommand {
	return CreateTagCommand{
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
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
