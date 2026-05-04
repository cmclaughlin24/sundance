package ports

import "github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"

type CreateSubmissionCommand struct {
	TenantID  string
	FormID    string
	VersionID string
	Payload   any
}

func NewCreateSubmssionCommand(tenantID, formID, versionID string, payload any) *CreateSubmissionCommand {
	return &CreateSubmissionCommand{
		TenantID:  tenantID,
		FormID:    formID,
		VersionID: versionID,
		Payload:   payload,
	}
}

type ReplaySubmissionCommand struct {
	TenantID    string
	ReferenceID domain.ReferenceID
}

func NewReplaySubmissionCommand(tenantID string, referenceID domain.ReferenceID) *ReplaySubmissionCommand {
	return &ReplaySubmissionCommand{
		TenantID:    tenantID,
		ReferenceID: referenceID,
	}
}
