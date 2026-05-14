package workers

import (
	"context"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/worker"
	"github.com/cmclaughlin24/sundance/backend/pkg/worker/elector"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type submissionJob struct {
	s       *domain.Submission
	service ports.SubmissionJobsService
}

func newSubmissionJob(service ports.SubmissionJobsService, s *domain.Submission) *submissionJob {
	return &submissionJob{
		s:       s,
		service: service,
	}
}

func (j *submissionJob) Process(ctx context.Context) error {
	return j.service.Process(ctx)
}

func NewDataSourcesBackgroundWorker(app *core.Application) (*worker.BackgroundWorker[*submissionJob], error) {
	bw, err := worker.NewBackgroundWorker[*submissionJob](
		worker.BgWithInterval[*submissionJob](1*time.Minute),
		worker.BgWithLogger[*submissionJob](app.Logger),
		worker.BgWithSize[*submissionJob](5),
		worker.BgWithFetchJobsFn[*submissionJob](newSubmissionWorkFn(app)),
		worker.BgWithElector[*submissionJob](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:submissions:elector"),
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
		dataSources, err := app.Services.SubmissionJobs.Find(ctx, &ports.FindSubmissionJobsQuery{})

		if err != nil {
			return nil, err
		}

		jobs := make([]*submissionJob, 0, len(dataSources))
		for _, ds := range dataSources {
			jobs = append(jobs, newSubmissionJob(
				app.Services.SubmissionJobs,
				ds,
			))
		}

		return jobs, nil
	}
}
