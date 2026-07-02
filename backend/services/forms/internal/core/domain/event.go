package domain

import (
	"encoding/json"
	"iter"
	"time"
)

// NOTE: Based on the Debezium Outbox Pattern (https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html)

type EventID string

type AggregateType string

type EventType string

type EventStatus string

const (
	EventStatusPending   EventStatus = "pending"
	EventStatusCompleted EventStatus = "completed"
	EventStatusError     EventStatus = "error"
)

type Event struct {
	ID            EventID
	AggregateID   string
	AggregateType AggregateType
	Type          EventType
	Status        EventStatus
	Payload       json.RawMessage
	Attempts      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewEvent(aggregateType AggregateType, aggregateID string, eventType EventType, payload json.RawMessage) Event {
	return Event{
		ID:            EventID(NewID()),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Type:          eventType,
		Status:        EventStatusPending,
		Payload:       payload,
		Attempts:      0,
		CreatedAt:     Now(),
		UpdatedAt:     Now(),
	}
}

func HydrateEvent(
	id EventID,
	aggregateType AggregateType,
	aggregateID string,
	eventType EventType,
	status EventStatus,
	payload json.RawMessage,
	attempts int,
	createdAt,
	updatedAt time.Time,
) Event {
	return Event{
		ID:            id,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Type:          eventType,
		Status:        status,
		Payload:       payload,
		Attempts:      attempts,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type withEvents struct {
	events []Event
}

func (we *withEvents) AddEvent(e Event) {
	we.events = append(we.events, e)
}

func (we *withEvents) DrainEvents() iter.Seq[Event] {
	events := we.events
	we.events = nil
	return func(yield func(Event) bool) {
		for _, e := range events {
			if !yield(e) {
				return
			}
		}
	}
}
