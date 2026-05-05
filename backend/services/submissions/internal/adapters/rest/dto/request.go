package dto

type SubmissionRequest struct {
	FormID    string `json:"formId"`
	VersionID string `json:"versionId"`
	Payload   any    `json:"payload"`
}
