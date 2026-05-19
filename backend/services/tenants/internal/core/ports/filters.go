package ports

import (
	"time"

	"sundance/backend/services/tenants/internal/core/domain"
)

type FindDataSourceJobsFilter struct {
	Types             []domain.DataSourceType
	Limit             int
	ExpiredAtOrBefore time.Time
}
