package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"sundance/backend/pkg/worker/elector"
)

var (
	ErrLoggerIsRequired      = errors.New("logger is required")
	ErrFetchJobsFnIsRequired = errors.New("fetchJobsFn is required")
)

type FetchJobsFn[J Job] func(context.Context) ([]J, error)

func NewBackgroundWorker[J Job](opts ...func(*BackgroundWorker[J])) (*BackgroundWorker[J], error) {
	bw := &BackgroundWorker[J]{
		elector:      elector.NewInMemoryElector(1 * time.Minute),
		interval:     1 * time.Minute,
		size:         5,
		failures:     0,
		failureLimit: nil,
	}

	for _, opt := range opts {
		opt(bw)
	}

	if err := bw.Validate(); err != nil {
		return nil, err
	}

	return bw, nil
}

func BgWithElector[J Job](elector elector.Elector) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.elector = elector
	}
}

func BgWithInterval[J Job](interval time.Duration) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.interval = interval
	}
}

func BgWithLogger[J Job](logger *slog.Logger) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.logger = logger
	}
}

func BgWithSize[J Job](size int) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.size = size
	}
}

func BgWithTimeout[J Job](timeout time.Duration) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.timeout = timeout
	}
}

func BgWithFetchJobsFn[J Job](fn FetchJobsFn[J]) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.fetchJobsFn = fn
	}
}

func BgWithFailureLimit[J Job](limit int) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.failureLimit = &limit
	}
}

type BackgroundWorker[J Job] struct {
	elector      elector.Elector
	interval     time.Duration
	logger       *slog.Logger
	size         int
	timeout      time.Duration
	fetchJobsFn  FetchJobsFn[J]
	failureLimit *int
	failures     int
}

func (bw *BackgroundWorker[J]) Start(ctx context.Context) {
	bw.logger.InfoContext(ctx, "background worker started", "pool_size", bw.size, "work_interval", bw.interval.String())

	ticker := time.NewTicker(bw.elector.GetInterval())
	defer ticker.Stop()

	isLeader := false
	var leaderCancel context.CancelFunc
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			bw.logger.InfoContext(ctx, "background worker stopping")

			if isLeader {
				_ = bw.elector.Release(context.Background())
			}

			if leaderCancel != nil {
				leaderCancel()
			}

			done := make(chan struct{})
			go func() { wg.Wait(); close(done) }()
			select {
			case <-done:
				bw.logger.InfoContext(ctx, "background worker gracefully shutdown")
			case <-time.After(30 * time.Second):
				bw.logger.ErrorContext(ctx, "background worker shutdown timed out")
			}

			return
		case <-ticker.C:
			if !isLeader {
				acquired, err := bw.elector.TryAcquire(ctx)

				if err != nil {
					bw.logger.ErrorContext(ctx, "failed to acquire leadership", "error", err)
					continue
				}

				if !acquired {
					continue
				}

				bw.logger.InfoContext(ctx, "leadership acquired")

				isLeader = true
				lctx, cancel := context.WithCancel(ctx)
				leaderCancel = cancel
				wg.Go(func() { bw.onLeader(lctx) })

				continue
			}

			if bw.shouldFailover() {
				isLeader = false
				bw.failures = 0

				if err := bw.elector.Release(ctx); err != nil {
					bw.logger.ErrorContext(ctx, "failed to release leadership", "error", err)
				}

				if leaderCancel != nil {
					leaderCancel()
					leaderCancel = nil
				}

				continue
			}

			ok, err := bw.elector.Renew(ctx)

			if err != nil || !ok {
				if err != nil {
					bw.logger.ErrorContext(ctx, "failed to renew leadership", "error", err)
				}

				bw.logger.WarnContext(ctx, "leadership lost")
				isLeader = false

				if leaderCancel != nil {
					leaderCancel()
					leaderCancel = nil
				}
			}
		}
	}
}

func (bw *BackgroundWorker[J]) onLeader(ctx context.Context) {
	ticker := time.NewTicker(bw.interval)
	defer ticker.Stop()

	pool := make(chan chan J)
	defer close(pool)

	for range bw.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](bw.logger), WorkerWithTimeout[J](bw.timeout))
		w.Start(ctx)
	}

	if err := bw.work(ctx, pool); err != nil {
		bw.recordFailure()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := bw.work(ctx, pool); err != nil && !errors.Is(err, context.Canceled) {
				bw.recordFailure()
			}
		}
	}
}

func (bw *BackgroundWorker[J]) work(ctx context.Context, pool chan chan J) error {
	jobs, err := bw.fetchJobsFn(ctx)

	if err != nil {
		bw.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)
		return err
	}

	bw.failures = 0
	bw.logger.DebugContext(ctx, "dispatching jobs", "count", len(jobs))

	for _, j := range jobs {
		select {
		case w := <-pool:
			w <- j
		case <-ctx.Done():
			return context.Canceled
		}
	}

	return nil
}

func (bw *BackgroundWorker[J]) Validate() error {
	if bw.logger == nil {
		return ErrLoggerIsRequired
	}

	if bw.fetchJobsFn == nil {
		return ErrFetchJobsFnIsRequired
	}

	return nil
}

func (bw *BackgroundWorker[J]) recordFailure() {
	bw.failures += 1
}

func (bw *BackgroundWorker[J]) shouldFailover() bool {
	if bw.failureLimit == nil {
		return false
	}
	return bw.failures >= *bw.failureLimit
}
