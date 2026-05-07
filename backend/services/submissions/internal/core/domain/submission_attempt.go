package domain

import "time"

type SubmissionAttemptID string

type SubmissionAttempt struct {
	ID           SubmissionAttemptID
	Attempt      int
	Result       string
	ErrorDetails any
	CreatedAt    time.Time
}

func NewSubmissionAttempt(
	attempt int,
	result string,
	errorDetails any,
) *SubmissionAttempt {
	return &SubmissionAttempt{
		ID:           SubmissionAttemptID(NewID()),
		Attempt:      attempt,
		Result:       result,
		ErrorDetails: errorDetails,
		CreatedAt:    Now(),
	}
}

func HydrateSubmissionAttempt(
	id SubmissionAttemptID,
	attempt int,
	result string,
	errorDetails any,
	createdAt time.Time,
) *SubmissionAttempt {
	return &SubmissionAttempt{
		ID:           id,
		Attempt:      attempt,
		Result:       result,
		ErrorDetails: errorDetails,
		CreatedAt:    createdAt,
	}
}
