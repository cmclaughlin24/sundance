package database

import "context"

type Database interface {
	Close() error
	BeginTx(context.Context) (context.Context, error)
	GetTx(context.Context) (any, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}

