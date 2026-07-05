package publishers

import (
	"context"
	"log/slog"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type InMemoryPublisher struct {
	logger *slog.Logger
}

func NewInMemoryPublisher(logger *slog.Logger) ports.Publisher {
	return &InMemoryPublisher{logger}
}

func (p *InMemoryPublisher) Publish(ctx context.Context, event domain.Event) error {
	p.logger.InfoContext(ctx, "in-memory publisher; event sunk (no-op)", "event_id", event.ID, "event_type", event.Type)

	return nil
}
