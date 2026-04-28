package stratreg

import "errors"

var (
	ErrStrategyNotFound = errors.New("strategy not found")
)

type StrategyRegistry[K comparable, U any] map[K]U

func New[K comparable, U any]() StrategyRegistry[K, U] {
	strategy := make(StrategyRegistry[K, U])
	return strategy
}

func (s StrategyRegistry[K, U]) Set(key K, strategy U) StrategyRegistry[K, U] {
	s[key] = strategy
	return s
}

func (s StrategyRegistry[K, U]) Get(key K) (U, error) {
	strategy, ok := s[key]

	if !ok {
		var z U
		return z, ErrStrategyNotFound
	}

	return strategy, nil
}
