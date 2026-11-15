package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type SubmissionDocument struct {
	ID            string                       `bson:"_id"`
	TenantID      string                       `bson:"tenant_id"`
	FormID        string                       `bson:"form_id"`
	VersionID     string                       `bson:"version_id"`
	ReferenceID   string                       `bson:"reference_id"`
	IdempotencyID string                       `bson:"idempotency_id"`
	Status        string                       `bson:"status"`
	CreatedAt     time.Time                    `bson:"created_at"`
	UpdatedAt     time.Time                    `bson:"updated_at"`
	Attempts      []*submissionAttemptDocument `bson:"attempts"`
	Facts         []*canonicalFact             `bson:"facts"`
	Sections      []*submissionSectionDocument `bson:"sections"`
	EvalContext   domain.EvaluationContext     `bson:"evaluation_context"`
}

func ToSubmissionDocument(s *domain.Submission) (*SubmissionDocument, error) {
	sections := make([]*submissionSectionDocument, 0, len(s.Sections))
	for _, sec := range s.Sections {
		sections = append(sections, toSubmissionSectionDocument(sec))
	}

	attempts := make([]*submissionAttemptDocument, 0, len(s.Attempts))
	for _, att := range s.Attempts {
		doc, err := toSubmissionAttemptDocument(att)

		if err != nil {
			return nil, err
		}

		attempts = append(attempts, doc)
	}

	facts := make([]*canonicalFact, 0, len(s.Facts))
	for _, f := range s.Facts {
		facts = append(facts, toCanonicalFactDocument(f))
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
		Facts:         facts,
		Sections:      sections,
		EvalContext:   s.GetEvalContext(),
	}, nil
}

func FromSubmissionDocument(s *SubmissionDocument) (*domain.Submission, error) {
	sections := make([]*domain.SubmissionSection, 0, len(s.Sections))
	for _, doc := range s.Sections {
		sections = append(sections, fromSubmissionSectionDocument(doc))
	}

	attempts := make([]*domain.SubmissionAttempt, 0, len(s.Attempts))
	for _, doc := range s.Attempts {
		attempts = append(attempts, fromSubmissionAttemptDocument(doc))
	}

	facts := make([]*domain.CanonicalFact, 0, len(s.Facts))
	for _, doc := range s.Facts {
		facts = append(facts, fromCanonicalFactDocument(doc))
	}

	return domain.HydrateSubmission(
		domain.SubmissionID(s.ID),
		s.TenantID,
		domain.FormID(s.FormID),
		domain.FormVersionID(s.VersionID),
		domain.ReferenceID(s.ReferenceID),
		domain.IdempotencyID(s.IdempotencyID),
		domain.SubmissionStatus(s.Status),
		sections,
		facts,
		attempts,
		s.EvalContext,
		s.CreatedAt,
		s.UpdatedAt,
	), nil
}
