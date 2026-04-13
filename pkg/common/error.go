package common

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrExists       = errors.New("already exists")
	ErrInvalidID    = errors.New("invalid id")
	ErrUnauthorized = errors.New("unauthorized access")
)
