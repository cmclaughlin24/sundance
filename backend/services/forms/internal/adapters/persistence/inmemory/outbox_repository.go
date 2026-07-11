package inmemory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type outboxEntry struct {
	event       *domain.Event
	lockedUntil time.Time
}

type inMemoryOutboxRepository struct {
	mu     sync.RWMutex
	events map[string]*outboxEntry
	logger *slog.Logger
}

func newInMemoryOutbox(logger *slog.Logger) ports.OutboxRepository {
	return &inMemoryOutboxRepository{
		events: make(map[string]*outboxEntry),
		logger: logger,
	}
}

func (r *inMemoryOutboxRepository) Claim(ctx context.Context, o ports.ClaimEventsOptions) ([]*domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	lockedUntil := now.Add(o.LeaseDuration)
	claimed := make([]*domain.Event, 0, o.BatchSize)

	for _, entry := range r.events {
		if len(claimed) >= o.BatchSize {
			break
		}

		e := entry.event

		if o.RetryLimit > 0 && e.Attempts >= o.RetryLimit {
			continue
		}

		if !o.CreatedAfter.IsZero() && e.CreatedAt.Before(o.CreatedAfter) {
			continue
		}

		eligible := e.Status == domain.EventStatusPending ||
			e.Status == domain.EventStatusError ||
			(e.Status == domain.EventStatusProcessing && entry.lockedUntil.Before(now))

		if !eligible {
			continue
		}

		e.Status = domain.EventStatusProcessing
		entry.lockedUntil = lockedUntil
		claimed = append(claimed, e)
	}

	return claimed, nil
}

func (r *inMemoryOutboxRepository) Upsert(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.events[string(event.ID)]
	if ok {
		existing.event = event
	} else {
		r.events[string(event.ID)] = &outboxEntry{event: event}
	}

	return event, nil
}
