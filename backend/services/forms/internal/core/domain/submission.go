package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"sundance/backend/pkg/common/validate"
)

type SubmissionID string

type ReferenceID string

type IdempotencyID string

type SubmissionStatus string

type FactMap map[string]any

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
	ErrDuplicateSection    = errors.New("each section response in a submission must be unique")
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
	Sections      []*SubmissionSection
	Facts         []*CanonicalFact
	Attempts      []*SubmissionAttempt
	evalContext   EvaluationContext
	withEvents
}

func NewSubmission(
	tenantID string,
	formID FormID,
	versionID FormVersionID,
	idempotencyID IdempotencyID,
	sections []*SubmissionSection,
	evalContext EvaluationContext,
) (*Submission, error) {
	s := &Submission{
		ID:            SubmissionID(NewID()),
		TenantID:      tenantID,
		FormID:        formID,
		VersionID:     versionID,
		ReferenceID:   ReferenceID(NewID()), // TODO: Investigate a more order number style implementation.
		IdempotencyID: idempotencyID,
		Status:        SubmissionStatusPending,
		Sections:      sections,
		Facts:         make([]*CanonicalFact, 0),
		Attempts:      make([]*SubmissionAttempt, 0),
		evalContext:   evalContext,
		CreatedAt:     Now(),
	}

	if err := s.validate(); err != nil {
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
	sections []*SubmissionSection,
	facts []*CanonicalFact,
	attempts []*SubmissionAttempt,
	evalContext EvaluationContext,
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
		Sections:      sections,
		Facts:         facts,
		Attempts:      attempts,
		evalContext:   evalContext,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (s *Submission) GetEvalContext() EvaluationContext {
	if s.evalContext == nil {
		return make(EvaluationContext)
	}
	return s.evalContext
}

func (s *Submission) GetSection(sectionID SectionID) *SubmissionSection {
	for _, section := range s.Sections {
		if section.SectionID == sectionID {
			return section.clone()
		}
	}

	return nil
}

func (s *Submission) Accept(facts []*CanonicalFact) {
	s.Status = SubmissionStatusAccepted
	s.Facts = facts
	s.UpdatedAt = Now()
	s.addAttempt(s.Status, nil)

	p, _ := json.Marshal(submissionAcceptedPayload{
		ReferenceID: s.ReferenceID,
		TenantID:    s.TenantID,
		FormID:      s.FormID,
		VersionID:   s.VersionID,
		Facts:       s.ToFactMap(),
	})
	s.addEvent(EventTypeSubmissionAccepted, p)
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

	p, _ := json.Marshal(submissionRejectedPayload{
		ReferenceID: s.ReferenceID,
		TenantID:    s.TenantID,
		FormID:      s.FormID,
		VersionID:   s.VersionID,
		Reason:      err.Error(),
	})
	s.addEvent(EventTypeSubmissionRejected, p)
}

func (s *Submission) Reset() {
	s.Status = SubmissionStatusPending
	s.UpdatedAt = Now()
}

func (s *Submission) ToFactMap() FactMap {
	return ToFactMap(s.Facts)
}

func (s *Submission) addAttempt(status SubmissionStatus, err error) {
	s.Attempts = append(s.Attempts, NewSubmissionAttempt(len(s.Attempts)+1, string(status), err))
}

func (s *Submission) addEvent(eventType EventType, payload json.RawMessage) {
	e := NewEvent(AggregateTypeSubmission, string(s.ID), eventType, payload)
	s.AddEvent(e)
}

func (s *Submission) validate() error {
	if err := validate.ValidateStruct(s); err != nil {
		return err
	}

	seen := make(map[SectionID]struct{})
	for _, section := range s.Sections {
		if _, ok := seen[section.SectionID]; ok {
			return ErrDuplicateSection
		}
		seen[section.SectionID] = struct{}{}
	}

	return nil
}

func ToFactMap(facts []*CanonicalFact) FactMap {
	result := make(FactMap)

	for _, fact := range facts {
		segments := strings.Split(fact.TagKey, pathSeparator)
		setNestedValue(result, segments, fact.Value, fact.CollectionIndex)
	}

	return result
}

func setNestedValue(node map[string]any, segments []string, value any, path []int) {
	raw := segments[0]
	isCollection := strings.HasSuffix(raw, TagCollectionSegment)
	key := strings.TrimSuffix(raw, TagCollectionSegment)

	if len(segments) == 1 {
		node[key] = value
		return
	}

	if isCollection {
		if len(path) == 0 {
			// NOTE: Path ends above this collection (missing value in a repeatable instance).
			if _, ok := node[key]; !ok {
				node[key] = nil
			}
			return
		}

		idx := path[0]
		slice, _ := node[key].([]map[string]any)
		for len(slice) <= idx {
			slice = append(slice, make(map[string]any))
		}

		node[key] = slice
		setNestedValue(slice[idx], segments[1:], value, path[1:])
		return
	}

	if _, ok := node[key]; !ok {
		node[key] = make(map[string]any)
	}

	setNestedValue(node[key].(map[string]any), segments[1:], value, path)
}

type submissionAcceptedPayload struct {
	ReferenceID ReferenceID   `json:"referenceId"`
	TenantID    string        `json:"tenantId"`
	FormID      FormID        `json:"formId"`
	VersionID   FormVersionID `json:"versionId"`
	Facts       FactMap       `json:"facts"`
}

type submissionFailedPayload struct {
	ReferenceID ReferenceID   `json:"referenceId"`
	TenantID    string        `json:"tenantId"`
	FormID      FormID        `json:"formId"`
	VersionID   FormVersionID `json:"versionId"`
	Reason      string        `json:"reason"`
}

type submissionRejectedPayload struct {
	ReferenceID ReferenceID   `json:"referenceId"`
	TenantID    string        `json:"tenantId"`
	FormID      FormID        `json:"formId"`
	VersionID   FormVersionID `json:"versionId"`
	Reason      string        `json:"reason"`
}
