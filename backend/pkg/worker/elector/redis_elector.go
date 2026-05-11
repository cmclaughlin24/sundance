package elector

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var renewScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

var releaseScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

type RedisElector struct {
	client     redis.Cmdable
	key        string
	instanceID string
	ttl        time.Duration
	interval   time.Duration
}

func NewRedisElector(opts ...func(*RedisElector)) *RedisElector {
	e := &RedisElector{
		instanceID: NewID(),
		interval:   1 * time.Minute,
		ttl:        2 * time.Minute,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func RedisElectorWithClient(client redis.Cmdable) func(*RedisElector) {
	return func(e *RedisElector) {
		e.client = client
	}
}

func RedisElectorWithKey(key string) func(*RedisElector) {
	return func(e *RedisElector) {
		e.key = key
	}
}

func RedisElectorWithInstanceID(instanceID string) func(*RedisElector) {
	return func(e *RedisElector) {
		e.instanceID = instanceID
	}
}

func RedisElectorWithInterval(interval time.Duration) func(*RedisElector) {
	return func(e *RedisElector) {
		e.interval = interval
	}
}

func RedisElectorWithTTL(ttl time.Duration) func(*RedisElector) {
	return func(e *RedisElector) {
		e.ttl = ttl
	}
}

func (e *RedisElector) GetInterval() time.Duration {
	return e.interval
}

func (e *RedisElector) TryAcquire(ctx context.Context) (bool, error) {
	return e.client.SetNX(ctx, e.key, e.instanceID, e.ttl).Result()
}

func (e *RedisElector) Renew(ctx context.Context) (bool, error) {
	result, err := renewScript.Run(
		ctx,
		e.client,
		[]string{e.key},
		e.instanceID,
		e.ttl.Milliseconds(),
	).Int()

	return result == 1, err
}

func (e *RedisElector) Release(ctx context.Context) error {
	_, err := releaseScript.Run(
		ctx,
		e.client,
		[]string{e.key},
		e.instanceID,
	).Result()

	return err
}
