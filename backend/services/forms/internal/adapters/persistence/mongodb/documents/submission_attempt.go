package documents

import (
	"time"

	"sundance/backend/services/forms/internal/core/domain"
)

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

type submissionValueDocument struct {
	ElementID       string `bson:"element_id"`
	Value           any    `bson:"value"`
	CollectionIndex *int   `bson:"collection_index,omitempty"`
}

func toSubmissionValueDocument(sv *domain.SubmissionValue) *submissionValueDocument {
	return &submissionValueDocument{
		ElementID:       string(sv.ElementID),
		Value:           sv.Value,
		CollectionIndex: sv.CollectionIndex,
	}
}

func fromSubmissionValueDocument(doc *submissionValueDocument) *domain.SubmissionValue {
	return domain.HydrateSubmissionValue(domain.ElementID(doc.ElementID), doc.Value, doc.CollectionIndex)
}
