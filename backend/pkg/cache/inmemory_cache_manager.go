package cache

import (
	"context"
	"encoding/json"
	"sync"
)

type InMemoryCacheManager struct {
	mu sync.RWMutex
	cache map[string]string
}

func NewInMemoryCacheManager() CacheManager {
	return &InMemoryCacheManager{
		cache: make(map[string]string),
	}
}

func (m *InMemoryCacheManager) Get(ctx context.Context, key string, data any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.cache[key]

	if !ok {
		return nil
	}

	if val == "" {
		return nil
	}

	return json.Unmarshal([]byte(val), &data)
}

func (m *InMemoryCacheManager) Set(ctx context.Context, key string, data any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	m.cache[key] = string(out[:])

	return nil
}

func (m *InMemoryCacheManager) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.cache[key]; !ok {
		return nil
	}

	delete(m.cache, key)

	return nil
}
