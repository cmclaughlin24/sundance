package worker

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type WorkerID string

type Job interface {
	Process(context.Context)
}

type Worker[J Job] struct {
	ID         WorkerID
	JobChannel chan J
	WorkerPool chan chan J
	logger     *slog.Logger
}

func NewWorker[J Job](pool chan chan J, logger *slog.Logger) Worker[J] {
	return Worker[J]{
		ID:         WorkerID(uuid.NewString()),
		JobChannel: make(chan J, 1),
		WorkerPool: pool,
		logger:     logger,
	}
}

func (w *Worker[J]) Start(ctx context.Context) {
	wctx := SetWorkerContext(ctx, string(w.ID))

	go func() {
		for {
			w.logger.DebugContext(wctx, "looking for work")

			select {
			case w.WorkerPool <- w.JobChannel:
			case <-wctx.Done():
				return
			}

			select {
			case job := <-w.JobChannel:
				job.Process(wctx)
			case <-wctx.Done():
				return
			}
		}
	}()
}
