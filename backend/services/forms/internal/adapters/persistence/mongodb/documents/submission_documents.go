package documents

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SubmissionDocument struct {
	ID            string                       `bson:"_id"`
	TenantID      string                       `bson:"tenant_id"`
	FormID        string                       `bson:"form_id"`
	VersionID     string                       `bson:"version_id"`
	ReferenceID   string                       `bson:"reference_id"`
	IdempotencyID string                       `bson:"idempotency_id"`
	Status        string                       `bson:"status"`
	Payload       bson.Raw                     `bson:"payload"`
	CreatedAt     time.Time                    `bson:"created_at"`
	UpdatedAt     time.Time                    `bson:"updated_at"`
	Attempts      []*submissionAttemptDocument `bson:"attempts"`
}

func ToSubmissionDocument(s *domain.Submission) (*SubmissionDocument, error) {
	payload, err := bson.Marshal(s.Payload)

	if err != nil {
		return nil, err
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
		FormID:        s.FormID,
		VersionID:     s.VersionID,
		ReferenceID:   string(s.ReferenceID),
		IdempotencyID: string(s.IdempotencyID),
		Status:        string(s.Status),
		Payload:       payload,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		Attempts:      attempts,
	}, nil
}

func FromSubmissionDocument(s *SubmissionDocument) (*domain.Submission, error) {
	payload, err := parsePayload(s.Payload)
	if err != nil {
		return nil, err
	}

	attempts := make([]*domain.SubmissionAttempt, 0, len(s.Attempts))
	for _, doc := range s.Attempts {
		attempts = append(attempts, fromSubmissionAttemptDocument(doc))
	}

	return domain.HydrateSubmission(
		domain.SubmissionID(s.ID),
		s.TenantID,
		s.FormID,
		s.VersionID,
		domain.ReferenceID(s.ReferenceID),
		domain.IdempotencyID(s.IdempotencyID),
		domain.SubmissionStatus(s.Status),
		payload,
		s.CreatedAt,
		s.UpdatedAt,
		attempts,
	), nil
}

type submissionAttemptDocument struct {
	ID           string    `bson:"_id"`
	Attempt      int       `bson:"attempt"`
	Result       string    `bson:"result"`
	ErrorDetails bson.Raw  `bson:"error_details"`
	CreatedAt    time.Time `bson:"created_at"`
}

func toSubmissionAttemptDocument(att *domain.SubmissionAttempt) (*submissionAttemptDocument, error) {
	details, err := bson.Marshal(att.ErrorDetails)

	if err != nil {
		return nil, err
	}

	return &submissionAttemptDocument{
		ID:           string(att.ID),
		Attempt:      att.Attempt,
		Result:       att.Result,
		ErrorDetails: details,
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

func parsePayload(raw bson.Raw) (map[string]any, error) {
	var payload map[string]any

	if err := bson.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}
