package ports

type CreateSubmissionCommand struct {
	TenantID  string
	FormID    string
	VersionID string
	Payload   any
}

type ReplaySubmissionCommand struct{}
