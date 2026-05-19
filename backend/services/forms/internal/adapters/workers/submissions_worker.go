package workers

import (
	"context"
	"time"

	"sundance/backend/pkg/worker"
	"sundance/backend/pkg/worker/elector"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type submissionJob struct {
	id      domain.SubmissionID
	service ports.SubmissionJobsService
}

func newSubmissionJob(service ports.SubmissionJobsService, id domain.SubmissionID) *submissionJob {
	return &submissionJob{
		id:      id,
		service: service,
	}
}

func (j *submissionJob) Process(ctx context.Context) error {
	return j.service.Process(ctx, j.id)
}

func NewSubmissionsBackgroundWorker(app *core.Application) (*worker.BackgroundWorker[*submissionJob], error) {
	bw, err := worker.NewBackgroundWorker[*submissionJob](
		worker.BgWithInterval[*submissionJob](1*time.Minute),
		worker.BgWithLogger[*submissionJob](app.Logger),
		worker.BgWithSize[*submissionJob](5),
		worker.BgWithFetchJobsFn[*submissionJob](newSubmissionWorkFn(app)),
		worker.BgWithElector[*submissionJob](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:forms:elector"),
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

func newSubmissionWorkFn(app *core.Application) worker.FetchJobsFn[*submissionJob] {
	return func(ctx context.Context) ([]*submissionJob, error) {
		ids, err := app.Services.SubmissionJobs.Find(ctx, ports.NewFindSubmissionJobsQuery(0))

		if err != nil {
			return nil, err
		}

		jobs := make([]*submissionJob, 0, len(ids))
		for _, id := range ids {
			jobs = append(jobs, newSubmissionJob(app.Services.SubmissionJobs, id))
		}

		return jobs, nil
	}
}
