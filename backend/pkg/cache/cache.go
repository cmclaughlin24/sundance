package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheType string

const (
	CacheTypeInMemory CacheType = "inmemory"
	CacheTypeRedis    CacheType = "redis"
)

type bootstrapFn func(CacheOptions, *slog.Logger) (CacheManager, error)

type CacheOptions any

type CacheSettings struct {
	Type    CacheType    `json:"type"`
	Options CacheOptions `json:"options"`
}

type CacheManager interface {
	Get(context.Context, string, any) error
	Set(context.Context, string, any) error
	Del(context.Context, string) error
	AcquireLock(context.Context, string, string, time.Duration) (bool, error)
	RenewLock(context.Context, string, string, time.Duration) (bool, error)
	ReleaseLock(context.Context, string, string) error
}

func Bootstrap(settings CacheSettings, logger *slog.Logger) (CacheManager, error) {
	var fn bootstrapFn

	switch settings.Type {
	case CacheTypeInMemory:
		fn = bootstrapInMemory
	case CacheTypeRedis:
		fn = bootstrapRedis
	}

	if fn == nil {
		return nil, fmt.Errorf("unknown cache type : %s", settings.Type)
	}

	return fn(settings.Options, logger)
}

func bootstrapInMemory(o CacheOptions, logger *slog.Logger) (CacheManager, error) {
	return NewInMemoryCacheManager(logger), nil
}

func bootstrapRedis(o CacheOptions, logger *slog.Logger) (CacheManager, error) {
	opts, err := parseOptions[RedisOptions](o)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", opts.Host, opts.Port),
	})

	if status := client.Ping(context.Background()); status.Err() != nil {
		return nil, status.Err()
	}

	return NewRedisCacheManager(client, logger), nil
}

func parseOptions[T CacheOptions](options CacheOptions) (T, error) {
	data, err := json.Marshal(options)

	if err != nil {
		return *new(T), err
	}

	var opts T

	if err := json.Unmarshal(data, &opts); err != nil {
		return *new(T), err
	}

	return opts, nil
}
