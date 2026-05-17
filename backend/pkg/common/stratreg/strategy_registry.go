package stratreg

import "errors"

var (
	ErrStrategyNotFound = errors.New("strategy not found")
)

type StrategyRegistry[K comparable, S any] map[K]S

func New[K comparable, S any]() StrategyRegistry[K, S] {
	strategy := make(StrategyRegistry[K, S])
	return strategy
}

func (s StrategyRegistry[K, S]) Set(key K, strategy S) StrategyRegistry[K, S] {
	s[key] = strategy
	return s
}

func (s StrategyRegistry[K, S]) Get(key K) (S, error) {
	strategy, ok := s[key]

	if !ok {
		var z S
		return z, ErrStrategyNotFound
	}

	return strategy, nil
}
