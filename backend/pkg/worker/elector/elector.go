package elector

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Creates a new random UUID and returns it as a string.
func NewID() string {
	return uuid.NewString()
}

type Elector interface {
	GetInterval() time.Duration
	TryAcquire(context.Context) (bool, error)
	Renew(context.Context) (bool, error)
	Release(context.Context) error
}
