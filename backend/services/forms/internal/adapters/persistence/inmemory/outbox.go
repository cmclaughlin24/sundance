package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryOutbox struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
	logger *slog.Logger
}

func newInMemoryOutbox(logger *slog.Logger) ports.Outbox {
	return &inMemoryOutbox{
		events: make(map[string]*domain.Event),
		logger: logger,
	}
}

func (r *inMemoryOutbox) Find(ctx context.Context) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*domain.Event, 0, len(r.events))

	for _, event := range r.events {
		events = append(events, event)
	}

	return events, nil
}

func (r *inMemoryOutbox) Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[string(event.ID)] = event

	return event, nil
}
