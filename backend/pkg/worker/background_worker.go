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

type BackgroundWorkerBuilder[J Job] struct {
	interval time.Duration
	logger   *slog.Logger
	size     int
	workFn   WorkFn[J]
	timeout  time.Duration
}

func NewBackgroundWorkerBuilder[J Job]() *BackgroundWorkerBuilder[J] {
	return &BackgroundWorkerBuilder[J]{}
}

func (b *BackgroundWorkerBuilder[J]) SetInterval(interval time.Duration) *BackgroundWorkerBuilder[J] {
	b.interval = interval
	return b
}

func (b *BackgroundWorkerBuilder[J]) SetLogger(logger *slog.Logger) *BackgroundWorkerBuilder[J] {
	b.logger = logger
	return b
}

func (b *BackgroundWorkerBuilder[J]) SetSize(size int) *BackgroundWorkerBuilder[J] {
	b.size = size
	return b
}

func (b *BackgroundWorkerBuilder[J]) SetWorkFn(fn WorkFn[J]) *BackgroundWorkerBuilder[J] {
	b.workFn = fn
	return b
}

func (b *BackgroundWorkerBuilder[J]) SetTimeout(timeout time.Duration) *BackgroundWorkerBuilder[J] {
	b.timeout = timeout
	return b
}

func (b *BackgroundWorkerBuilder[J]) Build() (*BackgroundWorker[J], error) {
	if b.logger == nil {
		return nil, ErrLoggerIsRequired
	}

	if b.workFn == nil {
		return nil, ErrWorkFnIsRequired
	}

	if b.interval == 0 {
		b.interval = 1 * time.Minute
	}

	if b.size == 0 {
		b.size = 5
	}

	return &BackgroundWorker[J]{
		interval: b.interval,
		logger:   b.logger,
		size:     b.size,
		workFn:   b.workFn,
		timeout:  b.timeout,
	}, nil
}

type BackgroundWorker[J Job] struct {
	interval time.Duration
	logger   *slog.Logger
	size     int
	workFn   WorkFn[J]
	timeout  time.Duration
}

func (wp *BackgroundWorker[J]) Start(ctx context.Context) {
	ticker := time.NewTicker(wp.interval)
	pool := make(chan chan J)

	defer close(pool)

	for range wp.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](wp.logger), WorkerWithTimeout[J](wp.timeout))
		w.Start(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := wp.workFn(ctx)

			if err != nil {
				continue
			}

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
