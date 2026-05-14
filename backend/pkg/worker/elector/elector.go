package elector

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Creates a new random UUID and returns it as a string.
func NewID() string {
	return uuid.NewString()
}

type ElectorType string

const (
	ElectorTypeInMemory ElectorType = "in-memory"
	ElectorTypeRedis    ElectorType = "redis"
)

type Elector interface {
	GetInterval() time.Duration
	TryAcquire(context.Context) (bool, error)
	Renew(context.Context) (bool, error)
	Release(context.Context) error
}

type CacheLocker interface {
	AcquireLock(context.Context, string, string, time.Duration) (bool, error)
	RenewLock(context.Context, string, string, time.Duration) (bool, error)
	ReleaseLock(context.Context, string, string) error
}
