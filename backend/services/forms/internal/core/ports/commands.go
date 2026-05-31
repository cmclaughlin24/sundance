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

func (c *DeleteCommand[T]) Validate() error {
	return validate.ValidateStruct(c)
}

type CreateCanonicalTagCommand struct {
	TenantID    string `validate:"required"`
	Key         string `validate:"required,nowhitespace"`
	DisplayName string `validate:"required"`
}

func NewCreateCanonicalTagCommand(tenantID, key, displayName string) CreateCanonicalTagCommand {
	return CreateCanonicalTagCommand{
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
	}
}

func (c CreateCanonicalTagCommand) Validate() error {
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

func NewCreateFormCommand(tenantID, name, description string) *CreateFormCommand {
	return &CreateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
	}
}

func (c *CreateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(tenantID string, id domain.FormID, name, description string) *UpdateFormCommand {
	return &UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

func (c *UpdateFormCommand) Validate() error {
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
	return validate.ValidateStruct(c)
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
	return validate.ValidateStruct(c)
}

type PublishFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewPublishFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) *PublishFormVersionCommand {
	return &PublishFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c *PublishFormVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RetireFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewRetireFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) *RetireFormVersionCommand {
	return &RetireFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c *RetireFormVersionCommand) Validate() error {
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
	return validate.ValidateStruct(c)
}

type ReplaySubmissionCommand struct {
	TenantID string              `validate:"required"`
	ID       domain.SubmissionID `validate:"required"`
}

func NewReplaySubmissionCommand(tenantID string, id domain.SubmissionID) *ReplaySubmissionCommand {
	return &ReplaySubmissionCommand{
		TenantID: tenantID,
		ID:       id,
	}
}

func (c *ReplaySubmissionCommand) Validate() error {
	return validate.ValidateStruct(c)
}
