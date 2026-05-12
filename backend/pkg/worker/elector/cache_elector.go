package elector

import (
	"context"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/cache"
)

type CacheElector struct {
	manager    cache.CacheManager
	key        string
	instanceID string
	ttl        time.Duration
	interval   time.Duration
}

func NewCacheElector(opts ...func(*CacheElector)) *CacheElector {
	e := &CacheElector{
		instanceID: NewID(),
		interval:   1 * time.Minute,
		ttl:        2 * time.Minute,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func CacheElectorWithManager(manager cache.CacheManager) func(*CacheElector) {
	return func(e *CacheElector) {
		e.manager = manager
	}
}

func CacheElectorWithKey(key string) func(*CacheElector) {
	return func(e *CacheElector) {
		e.key = key
	}
}

func CacheElectorWithInstanceID(instanceID string) func(*CacheElector) {
	return func(e *CacheElector) {
		e.instanceID = instanceID
	}
}

func CacheElectorWithInterval(interval time.Duration) func(*CacheElector) {
	return func(e *CacheElector) {
		e.interval = interval
	}
}

func CacheElectorWithTTL(ttl time.Duration) func(*CacheElector) {
	return func(e *CacheElector) {
		e.ttl = ttl
	}
}

func (e *CacheElector) GetInterval() time.Duration {
	return e.interval
}

func (e *CacheElector) TryAcquire(ctx context.Context) (bool, error) {
	return e.manager.AcquireLock(ctx, e.key, e.instanceID, e.ttl)
}

func (e *CacheElector) Renew(ctx context.Context) (bool, error) {
	return e.manager.RenewLock(ctx, e.key, e.instanceID, e.ttl)
}

func (e *CacheElector) Release(ctx context.Context) error {
	return e.manager.ReleaseLock(ctx, e.key, e.instanceID)
}
