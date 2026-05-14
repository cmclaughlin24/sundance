package domain

import "errors"

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
