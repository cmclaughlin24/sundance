package workers

import (
	"context"
	"sundance/backend/pkg/worker"
	"sundance/backend/pkg/worker/elector"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type outboxJob struct {
	event *domain.Event
}

func (j *outboxJob) Process(ctx context.Context) error {
	return nil
}

func newOutboxRelayBackgroundWorker(app *core.Application, opts ...func(*WorkerOptions)) (*worker.BackgroundWorker[*outboxJob], error) {
	options := newWorkerOptions(opts...)

	bw, err := worker.NewBackgroundWorker(
		worker.BgWithInterval[*outboxJob](time.Duration(options.Interval)*time.Minute),
		worker.BgWithLogger[*outboxJob](app.Logger),
		worker.BgWithSize[*outboxJob](options.PoolSize),
		worker.BgWithFetchJobsFn(newOutboxWorkFn(app, options.RetryLimit)),
		worker.BgWithElector[*outboxJob](elector.NewCacheElector(
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

func newOutboxWorkFn(app *core.Application, retryLimit int) worker.FetchJobsFn[*outboxJob] {
	return func(ctx context.Context) ([]*outboxJob, error) {
		return nil, nil
	}
}
