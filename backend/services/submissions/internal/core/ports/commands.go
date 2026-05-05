package ports

import "github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"

type CreateSubmissionCommand struct {
	TenantID  string
	FormID    string
	VersionID string
	Payload   any
}

func NewCreateSubmissionCommand(tenantID, formID, versionID string, payload any) *CreateSubmissionCommand {
	return &CreateSubmissionCommand{
		TenantID:  tenantID,
		FormID:    formID,
		VersionID: versionID,
		Payload:   payload,
	}
}

type ReplaySubmissionCommand struct {
	TenantID    string
	ReferenceID domain.SubmissionID
}
