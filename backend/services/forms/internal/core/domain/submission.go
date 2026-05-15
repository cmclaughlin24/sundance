package domain

import (
	"errors"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
)

type SubmissionID string

type ReferenceID string

type IdempotencyID string

type SubmissionStatus string

const (
	SubmissionStatusPending  SubmissionStatus = "pending"
	SubmissionStatusAccepted SubmissionStatus = "accepted"
	SubmissionStatusRejected SubmissionStatus = "rejected"
)

var (
	ErrDuplicateSubmission = errors.New("duplicate submissions")
)

type Submission struct {
	ID            SubmissionID
	TenantID      string    `validate:"required"`
	FormID        FormID    `validate:"required"`
	VersionID     VersionID `validate:"required"`
	ReferenceID   ReferenceID
	IdempotencyID IdempotencyID `validate:"required"`
	Status        SubmissionStatus
	Payload       map[string]any `validate:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Attempts      []*SubmissionAttempt
}

func NewSubmission(
	tenantID string,
	formID FormID,
	versionID VersionID,
	idempotencyID IdempotencyID,
	payload map[string]any,
) (*Submission, error) {
	s := &Submission{
		ID:            SubmissionID(NewID()),
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		ReferenceID:   ReferenceID(NewID()), // TODO: Investigate a more order number style implementation.
		IdempotencyID: idempotencyID,
		Status:        SubmissionStatusPending,
		Payload:       payload,
		CreatedAt:     Now(),
		Attempts:      make([]*SubmissionAttempt, 0),
	}

	if err := validate.ValidateStruct(s); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSubmission(
	id SubmissionID,
	tenantID string,
	formID FormID,
	versionID VersionID,
	referenceID ReferenceID,
	idempotencyID IdempotencyID,
	status SubmissionStatus,
	payload map[string]any,
	createdAt time.Time,
	updatedAt time.Time,
	attempts []*SubmissionAttempt,
) *Submission {
	return &Submission{
		ID:            id,
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		ReferenceID:   referenceID,
		IdempotencyID: idempotencyID,
		Status:        status,
		Payload:       payload,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Attempts:      attempts,
	}
}

func (s *Submission) Reset() {
	s.Status = SubmissionStatusPending
}
