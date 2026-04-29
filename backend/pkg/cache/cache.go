package cache

import (
	"context"
	"fmt"
)

type CacheType string

const (
	CacheTypeInMemory CacheType = "inmemory"
	CacheTypeRedis    CacheType = "redis"
)

type bootstrapFn func(CacheSettings) (CacheManager, error)

type CacheOptions any

type CacheSettings struct {
	Type    CacheType    `json:"type"`
	Options CacheOptions `json:"options"`
}

type CacheManager interface {
	Get(context.Context, string, any) error
	Set(context.Context, string, any) error
	Del(context.Context, string) error
}

func Bootstrap(settings CacheSettings) (CacheManager, error) {
	var fn bootstrapFn

	switch settings.Type {
	case CacheTypeInMemory:
		fn = bootstrapInMemory
	}

	if fn == nil {
		return nil, fmt.Errorf("unknown cache type : %s", settings.Type)
	}

	return fn(settings)
}

func bootstrapInMemory(_ CacheSettings) (CacheManager, error) {
	return NewInMemoryCacheManager(), nil
}
