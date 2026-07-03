package workers

import (
	"context"
	"log/slog"
	"sundance/backend/pkg/worker"
	"sundance/backend/pkg/worker/elector"
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
	if err := j.publisher.Publish(ctx); err != nil {
		j.event.Error()
	} else {
		j.event.Complete()
	}

	if _, err := j.outbox.Upsert(ctx, j.event); err != nil {
		return err
	}

	return nil
}

func newOutboxRelayBackgroundWorker(app *core.Application, opts ...func(*WorkerOptions)) (*worker.BackgroundWorker[*outboxMessage], error) {
	options := newWorkerOptions(opts...)

	bw, err := worker.NewBackgroundWorker(
		worker.BgWithInterval[*outboxMessage](time.Duration(options.Interval)*time.Minute),
		worker.BgWithLogger[*outboxMessage](app.Logger),
		worker.BgWithSize[*outboxMessage](options.PoolSize),
		worker.BgWithFetchJobsFn(newOutboxWorkFn(app, options.RetryLimit)),
		worker.BgWithElector[*outboxMessage](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:forms:elector:outbox"),
			elector.CacheElectorWithLocker(app.Cache),
			elector.CacheElectorWithInterval(1*time.Minute),
			elector.CacheElectorWithTTL(2*time.Minute),
		)),
	)

	if err != nil {
		return nil, err
	}

	return bw, nil
}

func newOutboxWorkFn(app *core.Application, retryLimit int) worker.FetchJobsFn[*outboxMessage] {
	return func(ctx context.Context) ([]*outboxMessage, error) {
		outbox := app.Outbox()

		app.Logger.DebugContext(ctx, "listing outbox messages")

		events, err := outbox.Find(ctx)
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
