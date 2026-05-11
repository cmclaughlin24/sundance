package elector

import (
	"context"
	"time"
)

type Elector interface {
	GetInterval() time.Duration
	TryAcquire(context.Context) (bool, error)
	Renew(context.Context) (bool, error)
	Release(context.Context) error
}
