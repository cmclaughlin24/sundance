package inmemory

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryOutboxRepository struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
	logger *slog.Logger
}

func newInMemoryOutbox(logger *slog.Logger) ports.OutboxRepository {
	return &inMemoryOutboxRepository{
		events: make(map[string]*domain.Event),
		logger: logger,
	}
}

func (r *inMemoryOutboxRepository) Find(ctx context.Context, filter ports.FindEventsFilter) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*domain.Event, 0, len(r.events))

	for _, event := range r.events {
		if len(filter.Statuses) > 0 && !slices.Contains(filter.Statuses, event.Status) {
			continue
		}

		if filter.RetryLimit > 0 && event.Attempts >= filter.RetryLimit {
			continue
		}

		if !filter.CreatedAfter.IsZero() && event.CreatedAt.Before(filter.CreatedAfter) {
			continue
		}

		events = append(events, event)
	}

	if filter.Take > 0 && len(events) > filter.Take {
		events = events[:filter.Take]
	}

	return events, nil
}

func (r *inMemoryOutboxRepository) Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[string(event.ID)] = event

	return event, nil
}
