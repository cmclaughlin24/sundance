package mongodb

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type submissionDocument struct {
	ID          string                       `bson:"_id"`
	TenantID    string                       `bson:"tenant_id"`
	FormID      string                       `bson:"form_id"`
	VersionID   string                       `bson:"version_id"`
	ReferenceID string                       `bson:"reference_id"`
	Status      string                       `bson:"status"`
	Payload     bson.Raw                     `bson:"payload"`
	CreatedAt   time.Time                    `bson:"created_at"`
	UpdatedAt   time.Time                    `bson:"updated_at"`
	Attempts    []*submissionAttemptDocument `bson:"attempts"`
}

func fromSubmissionDocument(s *submissionDocument) *domain.Submission {
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
		domain.SubmissionStatus(s.Status),
		s.Payload,
		s.CreatedAt,
		s.UpdatedAt,
		attempts,
	)
}

type submissionAttemptDocument struct {
	ID            string    `bson:"_id"`
	IdempotencyID string    `bson:"idempotency_id"`
	Attempt       int       `bson:"attempt"`
	Result        string    `bson:"result"`
	ErrorDetails  bson.Raw  `bson:"error_details"`
	CreatedAt     time.Time `bson:"created_at"`
}

func fromSubmissionAttemptDocument(att *submissionAttemptDocument) *domain.SubmissionAttempt {
	return domain.HydrateSubmissionAttempt(
		domain.SubmissionAttemptID(att.ID),
		domain.IdempotencyID(att.IdempotencyID),
		att.Attempt,
		att.Result,
		att.ErrorDetails,
		att.CreatedAt,
	)
}
