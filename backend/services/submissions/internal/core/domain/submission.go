package domain

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/google/uuid"
)

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
	Attempts    []*SubmissionAttempt
}

func NewSubmission(
	tenantID string,
	formID string,
	versionID string,
	payload any,
) (*Submission, error) {
	s := &Submission{
		ID:          SubmissionID(uuid.NewString()),
		TenantID:    tenantID,
		FormID:      formID,
		VersionID:   versionID,
		ReferenceID: ReferenceID(uuid.NewString()), // TODO: Investigate a more order number style implementation.
		Status:      "",                            // TODO: Implement submission state machine.
		Payload:     payload,
		CreatedAt:   Now(),
		Attempts:    make([]*SubmissionAttempt, 0),
	}

	if err := validate.ValidateStruct(s); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSubmission(
	id SubmissionID,
	tenantID string,
	formID string,
	versionID string,
	referenceID ReferenceID,
	status SubmissionStatus,
	payload any,
	createdAt time.Time,
	updatedAt time.Time,
	attempts []*SubmissionAttempt,
) *Submission {
	return &Submission{
		ID:          id,
		TenantID:    tenantID,
		FormID:      formID,
		VersionID:   versionID,
		ReferenceID: referenceID,
		Status:      status,
		Payload:     payload,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Attempts:    attempts,
	}
}
