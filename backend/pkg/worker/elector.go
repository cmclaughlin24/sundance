package worker

import "context"

type Elector interface {
	TryAcquire(context.Context) (bool, error)
	Renew(context.Context) (bool, error)
	Release(context.Context) error
}

type InMemoryElector struct{}

func NewInMemoryElector() *InMemoryElector {
	return &InMemoryElector{}
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
