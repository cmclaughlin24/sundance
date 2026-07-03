package workers

import (
	"context"
	"sundance/backend/services/forms/internal/core"
)

type WorkerOptions struct {
	Interval   int `json:"interval" env:"INTERVAL"`
	PoolSize   int `json:"poolSize" env:"POOL_SIZE"`
	RetryLimit int `json:"retryLimit" env:"RETRY_LIMIT"`
}

func Bootstrap(app *core.Application, settings WorkerOptions) (func(context.Context), error) {
	sw, err := newSubmissionsBackgroundWorker(
		app,
		WithInterval(settings.Interval),
		WithPoolSize(settings.PoolSize),
		WithRetryLimit(settings.RetryLimit),
	)

	if err != nil {
		return nil, err
	}

	ow, err := newOutboxRelayBackgroundWorker(
		app,
		WithInterval(settings.Interval),
		WithPoolSize(settings.PoolSize),
		WithRetryLimit(settings.RetryLimit),
	)

	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		go sw.Start(ctx)
		go ow.Start(ctx)
	}, nil
}

func newWorkerOptions(opts ...func(*WorkerOptions)) *WorkerOptions {
	o := &WorkerOptions{
		Interval:   1,
		PoolSize:   5,
		RetryLimit: 5,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithInterval(interval int) func(*WorkerOptions) {
	return func(ws *WorkerOptions) {
		ws.Interval = interval
	}
}

func WithPoolSize(size int) func(*WorkerOptions) {
	return func(ws *WorkerOptions) {
		ws.PoolSize = size
	}
}

func WithRetryLimit(limit int) func(*WorkerOptions) {
	return func(ws *WorkerOptions) {
		ws.RetryLimit = limit
	}
}
