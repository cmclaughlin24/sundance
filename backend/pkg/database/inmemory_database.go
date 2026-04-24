package database

import "context"

type txKeyType string

const txKey txKeyType = "tx"

type InMemoryDatabase struct{}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{}
}

func (db *InMemoryDatabase) Close() error {
	return nil
}

func (db *InMemoryDatabase) BeginTx(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, txKey, struct{}{}), nil
}

func (db *InMemoryDatabase) CommitTx(ctx context.Context) error {
	if tx := ctx.Value(txKey); tx == nil {
		return nil
	}

	return nil
}

func (db *InMemoryDatabase) RollbackTx(ctx context.Context) error {
	if tx := ctx.Value(txKey); tx == nil {
		return nil
	}

	return nil
}
