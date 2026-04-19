package strategy

import "errors"

var (
	ErrStrategyNotFound = errors.New("strategy not found")
)

type Strategies[K comparable, U any] map[K]U

func NewStrategies[K comparable, U any]() Strategies[K, U] {
	strategy := make(Strategies[K, U])
	return strategy
}

func (s Strategies[K, U]) Set(key K, strategy U) Strategies[K, U] {
	s[key] = strategy
	return s
}

func (s Strategies[K, U]) Get(key K) (U, error) {
	strategy, ok := s[key]

	if !ok {
		var z U
		return z, ErrStrategyNotFound
	}

	return strategy, nil
}
