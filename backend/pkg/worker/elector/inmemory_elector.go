package elector

import (
	"context"
	"time"
)

type InMemoryElector struct {
	interval time.Duration
}

func NewInMemoryElector(interval time.Duration) Elector {
	return &InMemoryElector{interval: interval}
}

func (e InMemoryElector) GetInterval() time.Duration {
	return e.interval
}

func (e *InMemoryElector) TryAcquire(context.Context) (bool, error) {
	return true, nil
}

func (e *InMemoryElector) Renew(context.Context) (bool, error) {
	return true, nil
}

func (e *InMemoryElector) Release(context.Context) error {
	return nil
}
