package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type submissionDocument struct {
	ID          string                       `json:"_id"`
	TenantID    string                       `json:"tenant_id"`
	FormID      string                       `json:"form_id"`
	VersionID   string                       `json:"version_id"`
	ReferenceID string                       `json:"reference_id"`
	Status      string                       `json:"status"`
	Payload     bson.Raw                     `json:"payload"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
	Attempts    []*submissionAttemptDocument `json:"attempts"`
}

type submissionAttemptDocument struct {
	ID            string    `json:"_id"`
	IdempotencyID string    `json:"idempotency_id"`
	Attempt       int       `json:"attempt"`
	Result        string    `json:"result"`
	ErrorDetails  bson.Raw  `json:"error_details"`
	CreatedAt     time.Time `json:"created_at"`
}
