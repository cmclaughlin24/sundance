package inmemory

import "context"

type InMemoryDatabase struct{}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{}
}

func (db *InMemoryDatabase) Close() error {
	return nil
}

func (db *InMemoryDatabase) BeginTx(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (db *InMemoryDatabase) GetTx(ctx context.Context) (any, error) {
	return nil, nil
}

func (db *InMemoryDatabase) CommitTx(ctx context.Context) error {
	return nil
}

func (db *InMemoryDatabase) RollbackTx(context.Context) error {
	return nil
}
