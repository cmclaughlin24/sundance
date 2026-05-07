package ports

import "github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"

type CreateSubmissionCommand struct {
	TenantID      string
	FormID        string
	VersionID     string
	IdempotencyID domain.IdempotencyID
	Payload       any
}

func NewCreateSubmissionCommand(
	tenantID,
	formID,
	versionID string,
	idempotencyID domain.IdempotencyID,
	payload any,
) *CreateSubmissionCommand {
	return &CreateSubmissionCommand{
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		IdempotencyID: idempotencyID,
		Payload:       payload,
	}
}

type ReplaySubmissionCommand struct {
	TenantID string
	ID       domain.SubmissionID
}

func NewReplaySubmissionCommand(tenantID string, id domain.SubmissionID) *ReplaySubmissionCommand {
	return &ReplaySubmissionCommand{
		TenantID: tenantID,
		ID:       id,
	}
}
