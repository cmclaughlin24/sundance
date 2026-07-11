package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

func NewPeriodicWorker[J Job](opts ...func(*PeriodicWorker[J])) (*PeriodicWorker[J], error) {
	pw := &PeriodicWorker[J]{
		interval: 1 * time.Minute,
		size:     5,
	}

	for _, opt := range opts {
		opt(pw)
	}

	if err := pw.Validate(); err != nil {
		return nil, err
	}

	return pw, nil
}

func PeriodicWithInterval[J Job](interval time.Duration) func(*PeriodicWorker[J]) {
	return func(bw *PeriodicWorker[J]) {
		bw.interval = interval
	}
}

func PeriodicWithLogger[J Job](logger *slog.Logger) func(*PeriodicWorker[J]) {
	return func(bw *PeriodicWorker[J]) {
		bw.logger = logger
	}
}

func PeriodicWithSize[J Job](size int) func(*PeriodicWorker[J]) {
	return func(bw *PeriodicWorker[J]) {
		bw.size = size
	}
}

func PeriodicWithTimeout[J Job](timeout time.Duration) func(*PeriodicWorker[J]) {
	return func(bw *PeriodicWorker[J]) {
		bw.timeout = timeout
	}
}

func PeriodicWithFetchJobsFn[J Job](fn FetchJobsFn[J]) func(*PeriodicWorker[J]) {
	return func(bw *PeriodicWorker[J]) {
		bw.fetchJobsFn = fn
	}
}

type PeriodicWorker[J Job] struct {
	logger      *slog.Logger
	interval    time.Duration
	timeout     time.Duration
	size        int
	fetchJobsFn FetchJobsFn[J]
}

func (pw *PeriodicWorker[J]) Start(ctx context.Context) {
	pw.logger.InfoContext(ctx, "periodic worker started", "pool_size", pw.size, "work_interval", pw.interval.String())

	ticker := time.NewTicker(pw.interval)
	defer ticker.Stop()

	pool := make(chan chan J, pw.size)
	defer close(pool)

	for range pw.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](pw.logger), WorkerWithTimeout[J](pw.timeout))
		w.Start(ctx)
	}

	if err := pw.work(ctx, pool); err != nil {
		pw.logger.WarnContext(ctx, "work failed", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			pw.logger.InfoContext(ctx, "periodic worker stopping")
			return
		case <-ticker.C:
			if err := pw.work(ctx, pool); err != nil && !errors.Is(err, context.Canceled) {
				pw.logger.WarnContext(ctx, "work failed", "error", err)
			}
		}
	}
}

func (pw *PeriodicWorker[J]) work(ctx context.Context, pool chan chan J) error {
	jobs, err := pw.fetchJobsFn(ctx)

	if err != nil {
		pw.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)
		return err
	}

	pw.logger.DebugContext(ctx, "dispatching jobs", "count", len(jobs))

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

func (pw *PeriodicWorker[J]) Validate() error {
	if pw.logger == nil {
		return ErrLoggerIsRequired
	}

	if pw.fetchJobsFn == nil {
		return ErrFetchJobsFnIsRequired
	}

	return nil
}
