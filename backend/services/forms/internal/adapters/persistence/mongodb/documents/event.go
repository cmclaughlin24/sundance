package documents

import (
	"encoding/json"
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type EventDocument struct {
	ID            string    `bson:"_id"`
	AggregateID   string    `bson:"aggregate_id"`
	AggregateType string    `bson:"aggregate_type"`
	Type          string    `bson:"type"`
	Status        string    `bson:"status"`
	Payload       string    `bson:"payload"`
	Attempts      int       `bson:"attempts"`
	LastError     *string   `bson:"error,omitempty"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

func ToEventDocument(e *domain.Event) *EventDocument {
	return &EventDocument{
		ID:            string(e.ID),
		AggregateID:   e.AggregateID,
		AggregateType: string(e.AggregateType),
		Type:          string(e.Type),
		Status:        string(e.Status),
		Payload:       string(e.Payload),
		Attempts:      e.Attempts,
		LastError:     e.LastError,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func FromEventDocument(doc *EventDocument) *domain.Event {
	e := domain.HydrateEvent(
		domain.EventID(doc.ID),
		domain.AggregateType(doc.AggregateType),
		doc.AggregateID,
		domain.EventType(doc.Type),
		domain.EventStatus(doc.Status),
		json.RawMessage(doc.Payload),
		doc.Attempts,
		doc.LastError,
		doc.CreatedAt,
		doc.UpdatedAt,
	)

	return &e
}
