package domain

import "errors"

var (
	ErrDuplicatePosition = errors.New("duplicate position")
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
