package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

type InMemoryCacheManager struct {
	mu      sync.RWMutex
	cache   map[string]string
	locks   map[string]string
	lockExp map[string]time.Time
	logger  *slog.Logger
}

func NewInMemoryCacheManager(logger *slog.Logger) CacheManager {
	return &InMemoryCacheManager{
		cache:   make(map[string]string),
		locks:   make(map[string]string),
		lockExp: make(map[string]time.Time),
		logger:  logger,
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

func (m *InMemoryCacheManager) AcquireLock(_ context.Context, key string, value string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if owner, held := m.locks[key]; held {
		if exp, ok := m.lockExp[key]; ok && time.Now().Before(exp) && owner != value {
			return false, nil
		}
	}

	m.locks[key] = value
	m.lockExp[key] = time.Now().Add(ttl)

	return true, nil
}

func (m *InMemoryCacheManager) ReleaseLock(_ context.Context, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if owner, held := m.locks[key]; held && owner == value {
		delete(m.locks, key)
		delete(m.lockExp, key)
	}

	return nil
}

func (m *InMemoryCacheManager) RenewLock(_ context.Context, key string, value string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	owner, held := m.locks[key]
	if !held || owner != value {
		return false, nil
	}

	if exp, ok := m.lockExp[key]; ok && time.Now().After(exp) {
		return false, nil
	}

	m.lockExp[key] = time.Now().Add(ttl)

	return true, nil
}
