package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type CreateSubmissionCommand struct {
	TenantID      string                    `validate:"required"`
	FormID        domain.FormID             `validate:"required"`
	VersionID     domain.FormVersionID      `validate:"required"`
	IdempotencyID domain.IdempotencyID      `validate:"required"`
	Values        []*domain.SubmissionValue `validate:"required,min=1"`
}

func NewCreateSubmissionCommand(
	tenantID string,
	formID domain.FormID,
	versionID domain.FormVersionID,
	idempotencyID domain.IdempotencyID,
	values []*domain.SubmissionValue,
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

type NormalizeSubmissionCommand struct {
	TenantID  string                    `validate:"required"`
	FormID    domain.FormID             `validate:"required"`
	VersionID domain.FormVersionID      `validate:"required"`
	Values    []*domain.SubmissionValue `validate:"required,min=1"`
}

func NewNormalizeSubmissionCommand(
	tenantID string,
	formID domain.FormID,
	versionID domain.FormVersionID,
	values []*domain.SubmissionValue,
) *NormalizeSubmissionCommand {
	return &NormalizeSubmissionCommand{
		TenantID:  tenantID,
		FormID:    formID,
		VersionID: versionID,
		Values:    values,
	}
}

func (c *NormalizeSubmissionCommand) Validate() error {
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
