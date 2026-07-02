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
	SubmissionStatusFailed   SubmissionStatus = "failed"

	AggregateTypeSubmission     AggregateType = "submission"
	EventTypeSubmissionAccepted EventType     = "accepted"
	EventTypeSubmissionRejected EventType     = "rejected"
)

var (
	ErrDuplicateSubmission = errors.New("duplicate submissions")
)

type Submission struct {
	ID            SubmissionID
	TenantID      string        `validate:"required"`
	FormID        FormID        `validate:"required"`
	VersionID     FormVersionID `validate:"required"`
	ReferenceID   ReferenceID
	IdempotencyID IdempotencyID `validate:"required"`
	Status        SubmissionStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Values        []*SubmissionFieldValue
	Facts         []*CanonicalFact
	Attempts      []*SubmissionAttempt
	withEvents
}

func NewSubmission(
	tenantID string,
	formID FormID,
	versionID FormVersionID,
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
		Facts:         make([]*CanonicalFact, 0),
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
	versionID FormVersionID,
	referenceID ReferenceID,
	idempotencyID IdempotencyID,
	status SubmissionStatus,
	values []*SubmissionFieldValue,
	facts []*CanonicalFact,
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
		Facts:         facts,
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

func (s *Submission) Accept(facts []*CanonicalFact) {
	s.Status = SubmissionStatusAccepted
	s.Facts = facts
	s.UpdatedAt = Now()
	s.addAttempt(s.Status, nil)
	// FIXME: Add domain event that the submission was accepted.
}

func (s *Submission) Fail(err error) {
	s.Status = SubmissionStatusFailed
	s.UpdatedAt = Now()
	s.addAttempt(s.Status, err)
}

func (s *Submission) Reject(err error) {
	s.Status = SubmissionStatusRejected
	s.UpdatedAt = Now()
	s.addAttempt(s.Status, err)
	// FIXME: Add domain event that the submission was rejected.
}

func (s *Submission) Reset() {
	s.Status = SubmissionStatusPending
	s.UpdatedAt = Now()
}

func (s *Submission) addAttempt(status SubmissionStatus, err error) {
	s.Attempts = append(s.Attempts, NewSubmissionAttempt(len(s.Attempts)+1, string(status), err))
}

func (s *Submission) ToFactMap() map[string]any {
	return nil
}

type SubmissionFieldValue struct {
	FieldID         FieldID
	Value           any
	CollectionIndex *int
}

func NewSubmissionFieldValue(fieldID FieldID, value any, collectionIndex *int) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID:         fieldID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}

func HydrateSubmissionFieldValue(fieldID FieldID, value any, collectionIndex *int) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID:         fieldID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}
