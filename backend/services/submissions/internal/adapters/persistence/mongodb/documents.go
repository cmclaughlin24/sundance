package mongodb

import (
	"time"

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

type submissionAttemptDocument struct {
	ID            string    `bson:"_id"`
	IdempotencyID string    `bson:"idempotency_id"`
	Attempt       int       `bson:"attempt"`
	Result        string    `bson:"result"`
	ErrorDetails  bson.Raw  `bson:"error_details"`
	CreatedAt     time.Time `bson:"created_at"`
}
