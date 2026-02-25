package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxUnitOfWorkFactory struct {
	pool *pgxpool.Pool
}

func NewUWFactory(ctx context.Context, pool *pgxpool.Pool) *pgxUnitOfWorkFactory {
	return &pgxUnitOfWorkFactory{pool: pool}
}

func (f *pgxUnitOfWorkFactory) Exec(ctx context.Context) (NoTxUnitOfWork, error) {
	return &pgxUnitOfWorkNoTx{pool: f.pool}, nil
}

func (f *pgxUnitOfWorkFactory) Begin(ctx context.Context) (TxUnitOfWork, error) {
	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &pgxUnitOfWorkTX{tx: tx}, nil
}
