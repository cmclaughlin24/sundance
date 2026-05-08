package dto

type SubmissionRequest struct {
	FormID    string `json:"formId" validate:"required,uuidv7"`
	VersionID string `json:"versionId" validate:"required,uuidv7"`
	Payload   any    `json:"payload" validate:"required"`
}
