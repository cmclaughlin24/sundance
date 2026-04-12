package common

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exits")
	ErrInvalidID = errors.New("invalid id")
)
