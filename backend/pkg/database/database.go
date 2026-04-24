package database

import "context"

type Database interface {
	Close() error
	BeginTx(context.Context) (context.Context, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}
