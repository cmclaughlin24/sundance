package documents

import (
	"time"

	"sundance/backend/services/forms/internal/core/domain"
)

type SubmissionDocument struct {
	ID            string                          `bson:"_id"`
	TenantID      string                          `bson:"tenant_id"`
	FormID        string                          `bson:"form_id"`
	VersionID     string                          `bson:"version_id"`
	ReferenceID   string                          `bson:"reference_id"`
	IdempotencyID string                          `bson:"idempotency_id"`
	Status        string                          `bson:"status"`
	CreatedAt     time.Time                       `bson:"created_at"`
	UpdatedAt     time.Time                       `bson:"updated_at"`
	Attempts      []*submissionAttemptDocument    `bson:"attempts"`
	Values        []*submissionFieldValueDocument `bson:"values"`
}

func ToSubmissionDocument(s *domain.Submission) (*SubmissionDocument, error) {
	values := make([]*submissionFieldValueDocument, 0, len(s.Values))
	for _, doc := range s.Values {
		values = append(values, toSubmissionFieldValueDocument(doc))
	}

	attempts := make([]*submissionAttemptDocument, 0, len(s.Attempts))
	for _, att := range s.Attempts {
		doc, err := toSubmissionAttemptDocument(att)

		if err != nil {
			return nil, err
		}

		attempts = append(attempts, doc)
	}

	return &SubmissionDocument{
		ID:            string(s.ID),
		TenantID:      s.TenantID,
		FormID:        string(s.FormID),
		VersionID:     string(s.VersionID),
		ReferenceID:   string(s.ReferenceID),
		IdempotencyID: string(s.IdempotencyID),
		Status:        string(s.Status),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		Attempts:      attempts,
		Values:        values,
	}, nil
}

func FromSubmissionDocument(s *SubmissionDocument) (*domain.Submission, error) {
	values := make([]*domain.SubmissionFieldValue, 0, len(s.Values))
	for _, doc := range s.Values {
		values = append(values, fromSubmissionFieldValueDocument(doc))
	}

	attempts := make([]*domain.SubmissionAttempt, 0, len(s.Attempts))
	for _, doc := range s.Attempts {
		attempts = append(attempts, fromSubmissionAttemptDocument(doc))
	}

	return domain.HydrateSubmission(
		domain.SubmissionID(s.ID),
		s.TenantID,
		domain.FormID(s.FormID),
		domain.VersionID(s.VersionID),
		domain.ReferenceID(s.ReferenceID),
		domain.IdempotencyID(s.IdempotencyID),
		domain.SubmissionStatus(s.Status),
		values,
		attempts,
		s.CreatedAt,
		s.UpdatedAt,
	), nil
}

type submissionAttemptDocument struct {
	ID           string    `bson:"_id"`
	Attempt      int       `bson:"attempt"`
	Result       string    `bson:"result"`
	ErrorDetails any       `bson:"error_details"`
	CreatedAt    time.Time `bson:"created_at"`
}

func toSubmissionAttemptDocument(att *domain.SubmissionAttempt) (*submissionAttemptDocument, error) {
	return &submissionAttemptDocument{
		ID:           string(att.ID),
		Attempt:      att.Attempt,
		Result:       att.Result,
		ErrorDetails: att.ErrorDetails,
		CreatedAt:    att.CreatedAt,
	}, nil
}

func fromSubmissionAttemptDocument(att *submissionAttemptDocument) *domain.SubmissionAttempt {
	return domain.HydrateSubmissionAttempt(
		domain.SubmissionAttemptID(att.ID),
		att.Attempt,
		att.Result,
		att.ErrorDetails,
		att.CreatedAt,
	)
}

type submissionFieldValueDocument struct {
	FieldID string `bson:"field_id"`
	Value   any    `bson:"value"`
}

func toSubmissionFieldValueDocument(fv *domain.SubmissionFieldValue) *submissionFieldValueDocument {
	return &submissionFieldValueDocument{
		FieldID: string(fv.FieldID),
		Value:   fv.Value,
	}
}

func fromSubmissionFieldValueDocument(doc *submissionFieldValueDocument) *domain.SubmissionFieldValue {
	return domain.HydrateSubmissionFieldValue(domain.FieldID(doc.FieldID), doc.Value)
}
