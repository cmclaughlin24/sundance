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
