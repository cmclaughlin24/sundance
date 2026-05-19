package domain

import (
	"errors"
	"slices"
	"time"

	"sundance/backend/pkg/common/validate"
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Values        []*SubmissionFieldValue
	Attempts      []*SubmissionAttempt
}

func NewSubmission(
	tenantID string,
	formID FormID,
	versionID VersionID,
	idempotencyID IdempotencyID,
	values []*SubmissionFieldValue,
) (*Submission, error) {
	s := &Submission{
		ID:            SubmissionID(NewID()),
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		ReferenceID:   ReferenceID(NewID()), // TODO: Investigate a more order number style implementation.
		IdempotencyID: idempotencyID,
		Status:        SubmissionStatusPending,
		Values:        values,
		Attempts:      make([]*SubmissionAttempt, 0),
		CreatedAt:     Now(),
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
	values []*SubmissionFieldValue,
	attempts []*SubmissionAttempt,
	createdAt time.Time,
	updatedAt time.Time,
) *Submission {
	return &Submission{
		ID:            id,
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		ReferenceID:   referenceID,
		IdempotencyID: idempotencyID,
		Status:        status,
		Values:        values,
		Attempts:      attempts,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (s *Submission) GetFieldValue(id FieldID) (*SubmissionFieldValue, bool) {
	idx := slices.IndexFunc(s.Values, func(fv *SubmissionFieldValue) bool {
		return id == fv.FieldID
	})

	if idx == -1 {
		return nil, false
	}

	return s.Values[idx], true
}

func (s *Submission) Reset() {
	s.Status = SubmissionStatusPending
}

type SubmissionFieldValue struct {
	FieldID FieldID
	Value   any
}

func NewSubmissionFieldValue(fieldID FieldID, value any) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID: fieldID,
		Value:   value,
	}
}

func HydrateSubmissionFieldValue(fieldID FieldID, value any) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID: fieldID,
		Value:   value,
	}
}
