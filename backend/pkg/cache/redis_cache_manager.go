package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

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

type RedisCacheManager struct {
	client redis.Cmdable
	logger *slog.Logger
}

func NewRedisCacheManager(client redis.Cmdable, logger *slog.Logger) CacheManager {
	return &RedisCacheManager{
		client: client,
		logger: logger,
	}
}

func (m *RedisCacheManager) Get(ctx context.Context, key string, data any) error {
	val, err := m.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			m.logger.DebugContext(ctx, "cache miss", "key", key)
			return ErrCacheMiss
		}

		m.logger.ErrorContext(ctx, "failed to get cache entry", "key", key, "error", err)
		return err
	}

	if val == "" {
		m.logger.DebugContext(ctx, "cache miss", "key", key)
		return ErrCacheMiss
	}

	if err := json.Unmarshal([]byte(val), &data); err != nil {
		m.logger.ErrorContext(ctx, "failed to unmarshal cache entry", "key", key, "error", err)
		return err
	}

	return nil
}

func (m *RedisCacheManager) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	out, err := json.Marshal(data)

	if err != nil {
		m.logger.ErrorContext(ctx, "failed to marshal cache entry", "key", key, "error", err)
		return err
	}

	status := m.client.Set(ctx, key, out, ttl)

	if status.Err() != nil {
		m.logger.ErrorContext(ctx, "failed to set cache entry", "key", key, "error", status.Err())
		return status.Err()
	}

	return nil
}

func (m *RedisCacheManager) Del(ctx context.Context, key string) error {
	cmd := m.client.Del(ctx, key)

	if cmd.Err() != nil {
		m.logger.ErrorContext(ctx, "failed to delete cache entry", "key", key, "error", cmd.Err())
	}

	return cmd.Err()
}

func (m *RedisCacheManager) AcquireLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	acquired, err := m.client.SetNX(ctx, key, value, ttl).Result()

	if err != nil {
		m.logger.ErrorContext(ctx, "failed to acquire lock", "key", key, "error", err)
	}

	return acquired, err
}

func (m *RedisCacheManager) RenewLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	result, err := renewScript.Run(
		ctx,
		m.client,
		[]string{key},
		value,
		ttl.Milliseconds(),
	).Int()

	if err != nil {
		m.logger.ErrorContext(ctx, "failed to renew lock", "key", key, "error", err)
	}

	return result == 1, err
}

func (m *RedisCacheManager) ReleaseLock(ctx context.Context, key string, value string) error {
	_, err := releaseScript.Run(
		ctx,
		m.client,
		[]string{key},
		value,
	).Result()

	if err != nil {
		m.logger.ErrorContext(ctx, "failed to release lock", "key", key, "error", err)
	}

	return err
}
