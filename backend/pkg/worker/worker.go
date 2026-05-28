package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var (
	// Package declaration for the current time function. Allows for easier testing by enabling the injection of a
	// mock time function.
	Now = time.Now

	// Creates a new random UUID and returns it as a string.
	NewID = uuid.NewString
)

type WorkerID string

type Job interface {
	Process(context.Context) error
}

type Worker[J Job] struct {
	ID         WorkerID
	JobChannel chan J
	WorkerPool chan chan J
	logger     *slog.Logger
	timeout    time.Duration
}

func NewWorker[J Job](opts ...func(*Worker[J])) Worker[J] {
	w := Worker[J]{
		ID:         WorkerID(NewID()),
		JobChannel: make(chan J, 1),
	}

	for _, opt := range opts {
		opt(&w)
	}

	return w
}

func WorkerWithPool[J Job](pool chan chan J) func(*Worker[J]) {
	return func(w *Worker[J]) {
		w.WorkerPool = pool
	}
}

func WorkerWithLogger[J Job](logger *slog.Logger) func(*Worker[J]) {
	return func(w *Worker[J]) {
		w.logger = logger
	}
}

func WorkerWithTimeout[J Job](timeout time.Duration) func(*Worker[J]) {
	return func(w *Worker[J]) {
		w.timeout = timeout
	}
}

func (w *Worker[J]) Start(ctx context.Context) {
	wctx := SetWorkerContext(ctx, string(w.ID))
	w.logger.DebugContext(wctx, "worker started")

	go func() {
		for {
			select {
			case w.WorkerPool <- w.JobChannel:
			case <-wctx.Done():
				w.logger.DebugContext(wctx, "worker stopping")
				return
			}

			select {
			case job := <-w.JobChannel:
				start := Now()
				jctx, cancel := w.setJobTimeout(wctx)

				if err := w.process(jctx, job); err != nil {
					w.logger.ErrorContext(jctx, "failed to process job", "error", err)
				}

				w.logger.InfoContext(jctx, "job processed", "duration", fmt.Sprintf("%d", time.Since(start).Milliseconds()))
				cancel()
			case <-wctx.Done():
				w.logger.DebugContext(wctx, "worker stopping")
				return
			}
		}
	}()
}

func (w *Worker[J]) process(ctx context.Context, job J) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.ErrorContext(ctx, "recovering from panic", "error", r)
		}
	}()

	return job.Process(ctx)
}

func (w *Worker[J]) setJobTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if w.timeout == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, w.timeout)
}
