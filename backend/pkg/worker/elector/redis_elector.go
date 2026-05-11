package elector

import (
	"context"
	"time"
)

type RedisElector struct {
	interval time.Duration
}

func NewRedisElector(interval time.Duration) *RedisElector {
	return &RedisElector{
		interval: interval,
	}
}

func (e RedisElector) GetInterval() time.Duration {
	return e.interval
}

func (e *RedisElector) TryAcquire(context.Context) (bool, error) {
	return true, nil
}

func (e *RedisElector) Renew(context.Context) (bool, error) {
	return true, nil
}

func (e *RedisElector) Release(context.Context) error {
	return nil
}
