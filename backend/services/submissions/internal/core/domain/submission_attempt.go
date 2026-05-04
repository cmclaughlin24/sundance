package domain

import "time"

type SubmissionAttemptID string

type IdempotencyID string

type SubmissionAttempt struct {
	ID            SubmissionAttemptID
	IdempotencyID IdempotencyID
	Attempt       int
	Result        string
	ErrorDetails  any
	CreatedAt     time.Time
}

func HydrateSubmissionAttempt(
	id SubmissionAttemptID,
	idempotencyID IdempotencyID,
	attempt int,
	result string,
	errorDetails any,
	createdAt time.Time,
) *SubmissionAttempt {
	return &SubmissionAttempt{
		ID:            id,
		IdempotencyID: idempotencyID,
		Attempt:       attempt,
		Result:        result,
		ErrorDetails:  errorDetails,
		CreatedAt:     createdAt,
	}
}
