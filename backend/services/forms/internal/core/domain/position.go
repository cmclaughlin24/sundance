package domain

import "errors"

var (
	ErrDuplicatePosition = errors.New("duplicate position")
	ErrInvalidPosition   = errors.New("invalid position; must be greater than or equal to 0")
)

type withPosition struct {
	position int
}

func (wp *withPosition) GetPosition() int {
	return wp.position
}

func (wp *withPosition) SetPosition(position int) {
	wp.position = position
}

func isValidPosition(position int) bool {
	return position >= 0
}
