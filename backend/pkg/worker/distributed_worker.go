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
)


func NewDistributedWorker[J Job](opts ...func(*DistributedWorker[J])) (*DistributedWorker[J], error) {
	dw := &DistributedWorker[J]{
		elector:      elector.NewInMemoryElector(1 * time.Minute),
		interval:     1 * time.Minute,
		size:         5,
		failures:     0,
		failureLimit: nil,
	}

	for _, opt := range opts {
		opt(dw)
	}

	if err := dw.Validate(); err != nil {
		return nil, err
	}

	return dw, nil
}

func DistributedWithElector[J Job](elector elector.Elector) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.elector = elector
	}
}

func DistributedWithInterval[J Job](interval time.Duration) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.interval = interval
	}
}

func DistributedWithLogger[J Job](logger *slog.Logger) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.logger = logger
	}
}

func DistributedWithSize[J Job](size int) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.size = size
	}
}

func DistributedWithTimeout[J Job](timeout time.Duration) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.timeout = timeout
	}
}

func DistributedWithFetchJobsFn[J Job](fn FetchJobsFn[J]) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.fetchJobsFn = fn
	}
}

func DistributedWithFailureLimit[J Job](limit int) func(*DistributedWorker[J]) {
	return func(bw *DistributedWorker[J]) {
		bw.failureLimit = &limit
	}
}

type DistributedWorker[J Job] struct {
	elector      elector.Elector
	interval     time.Duration
	logger       *slog.Logger
	size         int
	timeout      time.Duration
	fetchJobsFn  FetchJobsFn[J]
	failureLimit *int
	failures     int
}

func (dw *DistributedWorker[J]) Start(ctx context.Context) {
	dw.logger.InfoContext(ctx, "distributed worker started", "pool_size", dw.size, "work_interval", dw.interval.String())

	ticker := time.NewTicker(dw.elector.GetInterval())
	defer ticker.Stop()

	isLeader := false
	var leaderCancel context.CancelFunc
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			dw.logger.InfoContext(ctx, "distributed worker stopping")

			if isLeader {
				_ = dw.elector.Release(context.Background())
			}

			if leaderCancel != nil {
				leaderCancel()
			}

			done := make(chan struct{})
			go func() { wg.Wait(); close(done) }()
			select {
			case <-done:
				dw.logger.InfoContext(ctx, "distributed worker gracefully shutdown")
			case <-time.After(30 * time.Second):
				dw.logger.ErrorContext(ctx, "distributed worker shutdown timed out")
			}

			return
		case <-ticker.C:
			if !isLeader {
				acquired, err := dw.elector.TryAcquire(ctx)

				if err != nil {
					dw.logger.ErrorContext(ctx, "failed to acquire leadership", "error", err)
					continue
				}

				if !acquired {
					continue
				}

				dw.logger.InfoContext(ctx, "leadership acquired")

				isLeader = true
				lctx, cancel := context.WithCancel(ctx)
				leaderCancel = cancel
				wg.Go(func() { dw.onLeader(lctx) })

				continue
			}

			if dw.shouldFailover() {
				isLeader = false
				dw.failures = 0

				if err := dw.elector.Release(ctx); err != nil {
					dw.logger.ErrorContext(ctx, "failed to release leadership", "error", err)
				}

				if leaderCancel != nil {
					leaderCancel()
					leaderCancel = nil
				}

				continue
			}

			ok, err := dw.elector.Renew(ctx)

			if err != nil || !ok {
				if err != nil {
					dw.logger.ErrorContext(ctx, "failed to renew leadership", "error", err)
				}

				dw.logger.WarnContext(ctx, "leadership lost")
				isLeader = false

				if leaderCancel != nil {
					leaderCancel()
					leaderCancel = nil
				}
			}
		}
	}
}

func (dw *DistributedWorker[J]) onLeader(ctx context.Context) {
	ticker := time.NewTicker(dw.interval)
	defer ticker.Stop()

	pool := make(chan chan J, dw.size)
	defer close(pool)

	for range dw.size {
		w := NewWorker(WorkerWithPool(pool), WorkerWithLogger[J](dw.logger), WorkerWithTimeout[J](dw.timeout))
		w.Start(ctx)
	}

	if err := dw.work(ctx, pool); err != nil {
		dw.recordFailure()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := dw.work(ctx, pool); err != nil && !errors.Is(err, context.Canceled) {
				dw.recordFailure()
			}
		}
	}
}

func (dw *DistributedWorker[J]) work(ctx context.Context, pool chan chan J) error {
	jobs, err := dw.fetchJobsFn(ctx)

	if err != nil {
		dw.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)
		return err
	}

	dw.failures = 0
	dw.logger.DebugContext(ctx, "dispatching jobs", "count", len(jobs))

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

func (dw *DistributedWorker[J]) Validate() error {
	if dw.logger == nil {
		return ErrLoggerIsRequired
	}

	if dw.fetchJobsFn == nil {
		return ErrFetchJobsFnIsRequired
	}

	return nil
}

func (dw *DistributedWorker[J]) recordFailure() {
	dw.failures += 1
}

func (dw *DistributedWorker[J]) shouldFailover() bool {
	if dw.failureLimit == nil {
		return false
	}
	return dw.failures >= *dw.failureLimit
}
