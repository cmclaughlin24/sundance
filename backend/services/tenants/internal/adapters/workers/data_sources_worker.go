package workers

import (
	"context"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/worker"
	"github.com/cmclaughlin24/sundance/backend/pkg/worker/elector"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type dataSourceJob struct {
	ds      *domain.DataSource
	service ports.DataSourceJobsService
}

func newDataSourceJob(service ports.DataSourceJobsService, ds *domain.DataSource) *dataSourceJob {
	return &dataSourceJob{
		ds:      ds,
		service: service,
	}
}

func (j *dataSourceJob) Process(ctx context.Context) error {
	return j.service.Process(ctx, ports.NewProcessDataSourceJobCommand(j.ds))
}

func NewDataSourcesBackgroundWorker(app *core.Application) (*worker.BackgroundWorker[*dataSourceJob], error) {
	bw, err := worker.NewBackgroundWorker[*dataSourceJob](
		worker.BgWithInterval[*dataSourceJob](1*time.Minute),
		worker.BgWithLogger[*dataSourceJob](app.Logger),
		worker.BgWithSize[*dataSourceJob](5),
		worker.BgWithFetchJobsFn[*dataSourceJob](newDataSourceWorkFn(app)),
		worker.BgWithElector[*dataSourceJob](elector.NewCacheElector(
			elector.CacheElectorWithKey("service:tenants:elector"),
			elector.CacheElectorWithManager(app.Cache),
			elector.CacheElectorWithInterval(1*time.Minute),
			elector.CacheElectorWithTTL(2*time.Minute),
		)),
	)

	if err != nil {
		return nil, err
	}

	return bw, nil
}

func newDataSourceWorkFn(app *core.Application) worker.FetchJobsFn[*dataSourceJob] {
	return func(ctx context.Context) ([]*dataSourceJob, error) {
		dataSources, err := app.Services.DataSourceJobs.Find(ctx, ports.NewFindDataSourceJobsQuery(0))

		if err != nil {
			return nil, err
		}

		jobs := make([]*dataSourceJob, 0, len(dataSources))
		for _, ds := range dataSources {
			jobs = append(jobs, newDataSourceJob(
				app.Services.DataSourceJobs,
				ds,
			))
		}

		return jobs, nil
	}
}
