package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/worker"
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

func NewDataSourcesBackgroundWorker(app *core.Application) *worker.BackgroundWorker[*dataSourceJob] {
	w, err := worker.NewBackgroundWorkerBuilder[*dataSourceJob]().
		SetInterval(time.Minute * 1).
		SetLogger(app.Logger).
		SetSize(5).
		SetWorkFn(newDataSourceWorkFn(app)).
		Build()

	if err != nil {
		panic(fmt.Errorf("failed to bootstrap DataSourcesBackgroundWorker; %w", err))
	}

	return w
}

func newDataSourceWorkFn(app *core.Application) worker.WorkFn[*dataSourceJob] {
	return func(ctx context.Context) ([]*dataSourceJob, error) {
		dataSources, err := app.Services.DataSourceJobs.Find(ctx)

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
