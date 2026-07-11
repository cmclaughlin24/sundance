package worker

import (
	"context"
	"errors"
)

var (
	ErrLoggerIsRequired      = errors.New("logger is required")
	ErrFetchJobsFnIsRequired = errors.New("fetchJobsFn is required")
)

type FetchJobsFn[J Job] func(context.Context) ([]J, error)
