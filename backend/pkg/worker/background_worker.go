package worker

import (
	"context"
	"log/slog"
	"time"
)

type BackgroundWorker[J Job] struct {
	interval  time.Duration
	logger    *slog.Logger
	size      int
	workFn    WorkFn[J]
}

func (wp *BackgroundWorker[J]) Start(ctx context.Context) {
	ticker := time.NewTicker(wp.interval)
	pool := make(chan chan J)

	defer close(pool)

	for range wp.size {
		w := NewWorker(pool, wp.logger)
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
