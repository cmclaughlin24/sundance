package database

import "context"

type Database interface {
	Close(context.Context) error
	BeginTx(context.Context) (context.Context, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}
