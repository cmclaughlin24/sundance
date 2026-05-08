package worker

import (
	"context"
	"log/slog"
	"time"
)

type WorkFn[J Job] func(context.Context) ([]J, error)

type BackgroundWorkerBuilder[J Job] struct {
	interval time.Duration
	logger   *slog.Logger
	size     int
	workFn   WorkFn[J]
}

func NewBackgroundWorkerBuilder[J Job]() *BackgroundWorkerBuilder[J] {
	return &BackgroundWorkerBuilder[J]{}
}

func (w *BackgroundWorkerBuilder[J]) SetInterval(interval time.Duration) *BackgroundWorkerBuilder[J] {
	w.interval = interval
	return w
}

func (w *BackgroundWorkerBuilder[J]) SetLogger(logger *slog.Logger) *BackgroundWorkerBuilder[J] {
	w.logger = logger
	return w
}

func (w *BackgroundWorkerBuilder[J]) SetSize(size int) *BackgroundWorkerBuilder[J] {
	w.size = size
	return w
}

func (w *BackgroundWorkerBuilder[J]) SetWorkFn(fn WorkFn[J]) *BackgroundWorkerBuilder[J] {
	w.workFn = fn
	return w
}

func (w *BackgroundWorkerBuilder[J]) Build() *BackgroundWorker[J] {
	return &BackgroundWorker[J]{
		interval: w.interval,
		logger:   w.logger,
		size:     w.size,
		workFn:   w.workFn,
	}
}
