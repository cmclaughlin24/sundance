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

type submissionFieldValueDocument struct {
	FieldID         string `bson:"field_id"`
	Value           any    `bson:"value"`
	CollectionIndex *int   `bson:"collection_index,omitempty"`
}

func toSubmissionFieldValueDocument(fv *domain.SubmissionFieldValue) *submissionFieldValueDocument {
	return &submissionFieldValueDocument{
		FieldID:         string(fv.FieldID),
		Value:           fv.Value,
		CollectionIndex: fv.CollectionIndex,
	}
}

func fromSubmissionFieldValueDocument(doc *submissionFieldValueDocument) *domain.SubmissionFieldValue {
	return domain.HydrateSubmissionFieldValue(domain.FieldID(doc.FieldID), doc.Value, doc.CollectionIndex)
}
