package domain

import "time"

type SubmissionID string

type ReferenceID string

type SubmissionStatus string

type Submission struct {
	ID          SubmissionID
	TenantID    string
	FormID      string
	VersionID   string
	ReferenceID ReferenceID
	Status      SubmissionStatus
	Payload     any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
