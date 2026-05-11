package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/worker/elector"
)

var (
	ErrLoggerIsRequired      = errors.New("logger is required")
	ErrFetchJobsFnIsRequired = errors.New("fetchJobsFn is required")
)

type FetchJobsFn[J Job] func(context.Context) ([]J, error)

func NewBackgroundWorker[J Job](opts ...func(*BackgroundWorker[J])) (*BackgroundWorker[J], error) {
	bw := &BackgroundWorker[J]{
		elector:      elector.NewInMemoryElector(1 * time.Minute),
		workInterval: 1 * time.Minute,
		size:         5,
	}

	for _, opt := range opts {
		opt(bw)
	}

	if err := bw.Validate(); err != nil {
		return nil, err
	}

	return bw, nil
}

func WithElector[J Job](elector elector.Elector) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.elector = elector
	}
}

func WithWorkInterval[J Job](interval time.Duration) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.workInterval = interval
	}
}

func WithLogger[J Job](logger *slog.Logger) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.logger = logger
	}
}

func WithSize[J Job](size int) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.size = size
	}
}

func WithTimeout[J Job](timeout time.Duration) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.timeout = timeout
	}
}

func WithFetchJobsFn[J Job](fn FetchJobsFn[J]) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.fetchJobsFn = fn
	}
}

type BackgroundWorker[J Job] struct {
	elector      elector.Elector
	workInterval time.Duration
	logger       *slog.Logger
	size         int
	timeout      time.Duration
	fetchJobsFn  FetchJobsFn[J]
}

func (bw *BackgroundWorker[J]) Start(ctx context.Context) {
	bw.logger.InfoContext(ctx, "background worker started", "pool_size", bw.size, "work_interval", bw.workInterval.String())

	ticker := time.NewTicker(bw.elector.GetInterval())
	defer ticker.Stop()

	isLeader := false
	var leaderCancel context.CancelFunc

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
				go bw.onLeader(lctx)

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
	ticker := time.NewTicker(bw.workInterval)
	defer ticker.Stop()

	pool := make(chan chan J)
	defer close(pool)

	for range bw.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](bw.logger), WorkerWithTimeout[J](bw.timeout))
		w.Start(ctx)
	}

	bw.work(ctx, pool)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bw.work(ctx, pool)
		}
	}
}

func (bw *BackgroundWorker[J]) work(ctx context.Context, pool chan chan J) {
	jobs, err := bw.fetchJobsFn(ctx)

	if err != nil {
		bw.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)
		return
	}

	bw.logger.DebugContext(ctx, "dispatching jobs", "count", len(jobs))

	for _, j := range jobs {
		select {
		case w := <-pool:
			w <- j
		case <-ctx.Done():
			return
		}
	}
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
