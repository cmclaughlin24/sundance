package domain

import (
	"errors"
	"slices"
)

var (
	ErrDuplicatePosition = errors.New("duplicate position")
	ErrInvalidPosition   = errors.New("invalid position; must be greater than or equal to 0")
)

type withPosition struct {
	position float32
}

func (wp *withPosition) GetPosition() float32 {
	return wp.position
}

func (wp *withPosition) SetPosition(position float32) {
	wp.position = position
}

func isValidPosition(position float32) bool {
	return position >= 0
}

type PositionGetter interface {
	GetPosition() float32
}

type PositionElements[T PositionGetter] = []T

func hasUniqueElements[T PositionGetter](elements PositionElements[T]) bool {
	seen := make(map[float32]struct{}, len(elements))

	for _, element := range elements {
		position := element.GetPosition()

		if _, exists := seen[position]; exists {
			return false
		}

		seen[position] = struct{}{}
	}

	return true
}

func sortElements[T PositionGetter](elements PositionElements[T]) {
	slices.SortFunc(elements, func(a, b T) int {
		if a.GetPosition() < b.GetPosition() {
			return -1
		}

		if a.GetPosition() > b.GetPosition() {
			return 1
		}

		return 0
	})
}
