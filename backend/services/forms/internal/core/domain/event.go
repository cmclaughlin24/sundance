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
	LastError     *string
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
	lastError *string,
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
		LastError:     lastError,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (e *Event) Complete() {
	e.Status = EventStatusCompleted
	e.LastError = nil
	e.Attempts += 1
}

func (e *Event) Error(err string) {
	e.Status = EventStatusError
	e.LastError = &err
	e.Attempts += 1
}

type HasEvents interface {
	PeekEvents() iter.Seq[Event]
	DrainEvents()
}

type withEvents struct {
	events []Event
}

func (we *withEvents) AddEvent(e Event) {
	we.events = append(we.events, e)
}

func (we *withEvents) PeekEvents() iter.Seq[Event] {
	return func(yield func(Event) bool) {
		for _, e := range we.events {
			if !yield(e) {
				return
			}
		}
	}
}

func (we *withEvents) DrainEvents() {
	we.events = nil
}
