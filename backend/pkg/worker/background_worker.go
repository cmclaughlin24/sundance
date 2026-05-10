package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var (
	ErrLoggerIsRequired = errors.New("logger is required")
	ErrWorkFnIsRequired = errors.New("workFn is required")
)

type WorkFn[J Job] func(context.Context) ([]J, error)

func NewBackgroundWorker[J Job](opts ...func(*BackgroundWorker[J])) (*BackgroundWorker[J], error) {
	bw := &BackgroundWorker[J]{
		elector:  NewInMemoryElector(),
		interval: 1 * time.Minute,
		size:     5,
	}

	for _, opt := range opts {
		opt(bw)
	}

	if err := bw.Validate(); err != nil {
		return nil, err
	}

	return bw, nil
}

func WithElector[J Job](elector Elector) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.elector = elector
	}
}

func WithInterval[J Job](interval time.Duration) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.interval = interval
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

func WithWorkFn[J Job](fn WorkFn[J]) func(*BackgroundWorker[J]) {
	return func(bw *BackgroundWorker[J]) {
		bw.workFn = fn
	}
}

type BackgroundWorker[J Job] struct {
	elector  Elector
	interval time.Duration
	logger   *slog.Logger
	size     int
	timeout  time.Duration
	workFn   WorkFn[J]
}

func (bw *BackgroundWorker[J]) Start(ctx context.Context) {
	bw.logger.InfoContext(ctx, "background worker started", "pool_size", bw.size, "interval", bw.interval)

	ticker := time.NewTicker(bw.interval)
	defer ticker.Stop()

	pool := make(chan chan J)
	defer close(pool)

	for range bw.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](bw.logger), WorkerWithTimeout[J](bw.timeout))
		w.Start(ctx)
	}

	isLeader := false

	for {
		select {
		case <-ctx.Done():
			bw.logger.InfoContext(ctx, "background worker stopping")

			if isLeader {
				_ = bw.elector.Release(context.Background())
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
			} else {
				ok, err := bw.elector.Renew(ctx)

				if err != nil || !ok {
					if err != nil {
						bw.logger.ErrorContext(ctx, "failed to renew leadership", "error", err)
					}

					bw.logger.WarnContext(ctx, "leadership lost")
					isLeader = false
					continue
				}
			}

			jobs, err := bw.workFn(ctx)

			if err != nil {
				bw.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)
				continue
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
	}
}

func (bw *BackgroundWorker[J]) Validate() error {
	if bw.logger == nil {
		return ErrLoggerIsRequired
	}

	if bw.workFn == nil {
		return ErrWorkFnIsRequired
	}

	return nil
}
