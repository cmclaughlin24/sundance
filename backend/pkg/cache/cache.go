package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheType string

const (
	CacheTypeInMemory CacheType = "in-memory"
	CacheTypeRedis    CacheType = "redis"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type CacheCloser func() error

type CacheSettings struct {
	Type  CacheType     `json:"type" env:"TYPE"`
	Redis *RedisOptions `json:"redis" envPrefix:"REDIS_" env:",init"`
}

type CacheManager interface {
	Get(context.Context, string, any) error
	Set(context.Context, string, any, time.Duration) error
	Del(context.Context, string) error
}

func Bootstrap(settings CacheSettings, logger *slog.Logger) (CacheManager, CacheCloser, error) {
	switch settings.Type {
	case CacheTypeInMemory:
		return bootstrapInMemory(logger)
	case CacheTypeRedis:
		return bootstrapRedis(settings.Redis, logger)
	default:
		return nil, nil, fmt.Errorf("unknown cache type : %s", settings.Type)
	}

}

func bootstrapInMemory(logger *slog.Logger) (CacheManager, CacheCloser, error) {
	return NewInMemoryCacheManager(logger), func() error { return nil }, nil
}

func bootstrapRedis(options *RedisOptions, logger *slog.Logger) (CacheManager, CacheCloser, error) {
	if options == nil {
		return nil, nil, errors.New("redis options are required for redis cache driver")
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.Port),
	})

	if status := client.Ping(context.Background()); status.Err() != nil {
		return nil, nil, status.Err()
	}

	return NewRedisCacheManager(client, logger), client.Close, nil
}
