package workers

import (
	"context"
	"time"

	"sundance/backend/pkg/worker"
	"sundance/backend/pkg/worker/elector"
	"sundance/backend/services/tenants/internal/core"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"sundance/backend/services/tenants/internal/core/ports/commands"
)

type dataSourceJob struct {
	ds  *domain.DataSource
	api ports.DataSourceJobsAPI
}

func newDataSourceJob(api ports.DataSourceJobsAPI, ds *domain.DataSource) *dataSourceJob {
	return &dataSourceJob{
		ds:  ds,
		api: api,
	}
}

func (j *dataSourceJob) Process(ctx context.Context) error {
	return j.api.Process(ctx, commands.NewProcessDataSourceJobCommand(j.ds))
}

func newDataSourcesBackgroundWorker(app *core.Application, opts ...func(*WorkerOptions)) (*worker.BackgroundWorker[*dataSourceJob], error) {
	options := newWorkerOptions(opts...)

	bw, err := worker.NewBackgroundWorker(
		worker.BgWithInterval[*dataSourceJob](time.Duration(options.Interval)*time.Minute),
		worker.BgWithLogger[*dataSourceJob](app.Logger),
		worker.BgWithSize[*dataSourceJob](options.PoolSize),
		worker.BgWithFetchJobsFn(newDataSourceWorkFn(app, options.RetryLimit)),
		worker.BgWithElector[*dataSourceJob](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:tenants:elector"),
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

func newDataSourceWorkFn(app *core.Application, retryLimit int) worker.FetchJobsFn[*dataSourceJob] {
	return func(ctx context.Context) ([]*dataSourceJob, error) {
		dataSources, err := app.API.DataSourceJobs.Find(ctx, ports.NewFindDataSourceJobsQuery(0, retryLimit))

		if err != nil {
			return nil, err
		}

		jobs := make([]*dataSourceJob, 0, len(dataSources))
		for _, ds := range dataSources {
			jobs = append(jobs, newDataSourceJob(
				app.API.DataSourceJobs,
				ds,
			))
		}

		return jobs, nil
	}
}
