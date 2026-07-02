package domain

import (
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

type EventPayload map[string]any

type Event struct {
	ID            EventID
	AggregateID   string
	AggregateType AggregateType
	Type          EventType
	Status        EventStatus
	Payload       EventPayload
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewEvent(aggregateType AggregateType, aggregateID string, eventType EventType, payload EventPayload) Event {
	return Event{
		ID:            EventID(NewID()),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Type:          eventType,
		Status:        EventStatusPending,
		Payload:       payload,
		CreatedAt:     Now(),
		UpdatedAt:     Now(),
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
