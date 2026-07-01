package domain

import "iter"

type Event struct {
}

func NewEvent() Event {
	return Event{}
}

type withEvents struct {
	events []Event
}

func (we *withEvents) AddEvent(e Event) {
	we.events = append(we.events, e)
}

func (we *withEvents) DrainEvents() iter.Seq[Event] {
	return func(yield func(Event) bool) {
		for _, e := range we.events {
			if !yield(e) {
				return
			}
		}
	}
}
