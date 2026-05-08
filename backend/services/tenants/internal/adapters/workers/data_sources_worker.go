package workers

import (
	"context"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/worker"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type dataSourceJob struct {
	ds *domain.DataSource
}

func (j dataSourceJob) Process(context.Context) {
}

func NewDataSourcesBackgroundWorker(app *core.Application) *worker.BackgroundWorker[dataSourceJob] {
	return worker.NewBackgroundWorkerBuilder[dataSourceJob]().
		SetInterval(time.Second * 15).
		SetLogger(app.Logger).
		SetSize(5).
		SetWorkFn(func(ctx context.Context) ([]dataSourceJob, error) {
			return nil, nil
		}).
		Build()
}
