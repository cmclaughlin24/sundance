package elector

import (
	"context"
	"time"
)

type CacheElector struct {
	locker     CacheLocker
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

func CacheElectorWithLocker(locker CacheLocker) func(*CacheElector) {
	return func(e *CacheElector) {
		e.locker = locker
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
	return e.locker.AcquireLock(ctx, e.key, e.instanceID, e.ttl)
}

func (e *CacheElector) Renew(ctx context.Context) (bool, error) {
	return e.locker.RenewLock(ctx, e.key, e.instanceID, e.ttl)
}

func (e *CacheElector) Release(ctx context.Context) error {
	return e.locker.ReleaseLock(ctx, e.key, e.instanceID)
}
