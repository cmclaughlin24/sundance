package workers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/pkg/worker"
	"sundance/backend/pkg/worker/elector"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

type submissionJobOptions struct {
	retryLimit int
	backoff    time.Duration
}

type submissionJob struct {
	id      domain.SubmissionID
	api     ports.SubmissionJobsAPI
	logger  *slog.Logger
	options submissionJobOptions
}

func newSubmissionJob(api ports.SubmissionJobsAPI, logger *slog.Logger, id domain.SubmissionID, options submissionJobOptions) *submissionJob {
	return &submissionJob{
		id:      id,
		logger:  logger,
		api:     api,
		options: options,
	}
}

func (j *submissionJob) Process(ctx context.Context) error {
	backoff := j.options.backoff

	for attempt := 1; attempt <= j.options.retryLimit; attempt++ {
		err := j.api.Process(ctx, j.id)

		if err == nil {
			break
		}

		if !isRetryableError(err) {
			j.logger.WarnContext(ctx, "stopped due to non-retryable error", "submission_id", j.id, "attempt", attempt, "error", err)
			break
		}

		if attempt == j.options.retryLimit {
			j.logger.WarnContext(ctx, "retry limit reached", "submission_id", j.id, "attempt", attempt, "error", err)
			break
		}

		j.logger.WarnContext(ctx, "retrying submission job", "submission_id", j.id, "attempt", attempt, "error", err)
		time.Sleep(backoff)
		backoff *= 2
	}

	return nil
}

func newSubmissionsDistributedWorker(app *core.Application, opts ...func(*WorkerOptions)) (*worker.DistributedWorker[*submissionJob], error) {
	options := newWorkerOptions(opts...)

	bw, err := worker.NewDistributedWorker(
		worker.DistributedWithInterval[*submissionJob](time.Duration(options.Interval)*time.Minute),
		worker.DistributedWithLogger[*submissionJob](app.Logger),
		worker.DistributedWithSize[*submissionJob](options.PoolSize),
		worker.DistributedWithFetchJobsFn(newSubmissionWorkFn(app, submissionJobOptions{
			retryLimit: options.RetryLimit,
			backoff:    time.Second,
		})),
		worker.DistributedWithElector[*submissionJob](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:forms:elector:submissions"),
			elector.CacheElectorWithLocker(app.Cache),
			elector.CacheElectorWithInterval(1*time.Minute),
			elector.CacheElectorWithTTL(2*time.Minute),
		)),
	)

	if err != nil {
		return nil, err
	}

	return bw, nil
}

func newSubmissionWorkFn(app *core.Application, options submissionJobOptions) worker.FetchJobsFn[*submissionJob] {
	return func(ctx context.Context) ([]*submissionJob, error) {
		ids, err := app.API.SubmissionJobs.Find(ctx, ports.NewFindSubmissionJobsQuery(0))

		if err != nil {
			return nil, err
		}

		jobs := make([]*submissionJob, 0, len(ids))
		for _, id := range ids {
			jobs = append(jobs, newSubmissionJob(
				app.API.SubmissionJobs,
				app.Logger,
				id,
				options,
			))
		}

		return jobs, nil
	}
}

func isRetryableError(err error) bool {
	return !errors.Is(err, strategies.ErrElementValidation) &&
		!errors.Is(err, strategies.ErrElementRequired) &&
		!errors.Is(err, strategies.ErrElementTypeValue) &&
		!errors.Is(err, domain.ErrInvalidVersionStatus) &&
		!errors.Is(err, stratreg.ErrStrategyNotFound)
}
