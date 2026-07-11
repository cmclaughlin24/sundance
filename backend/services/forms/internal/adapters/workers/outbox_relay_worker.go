package workers

import (
	"context"
	"log/slog"
	"sundance/backend/pkg/worker"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"time"
)

type outboxMessage struct {
	event     *domain.Event
	logger    *slog.Logger
	outbox    ports.OutboxRepository
	publisher ports.Publisher
}

func (j *outboxMessage) Process(ctx context.Context) error {
	if err := j.publisher.Publish(ctx, *j.event); err != nil {
		j.event.Error(err.Error())
	} else {
		j.event.Complete()
	}

	if _, err := j.outbox.Upsert(ctx, j.event); err != nil {
		j.logger.ErrorContext(ctx, "failed to persist outbox event", "event_id", j.event.ID, "attempt", j.event.Attempts, "error", err)
		return err
	}

	return nil
}

func newOutboxRelayPeriodicWorker(app *core.Application, opts ...func(*WorkerOptions)) (*worker.PeriodicWorker[*outboxMessage], error) {
	options := newWorkerOptions(opts...)

	pw, err := worker.NewPeriodicWorker(
		worker.PeriodicWithInterval[*outboxMessage](time.Duration(options.Interval)*time.Minute),
		worker.PeriodicWithLogger[*outboxMessage](app.Logger),
		worker.PeriodicWithSize[*outboxMessage](options.PoolSize),
		worker.PeriodicWithFetchJobsFn(newOutboxWorkFn(app, 10, options.RetryLimit)),
	)

	if err != nil {
		return nil, err
	}

	return pw, nil
}

func newOutboxWorkFn(app *core.Application, batchSize, retryLimit int) worker.FetchJobsFn[*outboxMessage] {
	return func(ctx context.Context) ([]*outboxMessage, error) {
		outbox := app.Outbox()

		app.Logger.DebugContext(ctx, "listing outbox messages")

		events, err := outbox.Claim(ctx, ports.ClaimEventsOptions{
			RetryLimit:    retryLimit,
			CreatedAfter:  time.Now().Add(-24 * time.Hour),
			BatchSize:     batchSize,
			LeaseDuration: 5 * time.Minute,
		})
		if err != nil {
			app.Logger.ErrorContext(ctx, "failed to retrieve outbox messages", "error", err)
			return nil, err
		}

		jobs := make([]*outboxMessage, 0, len(events))
		for _, e := range events {
			jobs = append(jobs, &outboxMessage{
				event:     e,
				logger:    app.Logger,
				outbox:    outbox,
				publisher: app.Publisher,
			})
		}

		return jobs, nil
	}
}
