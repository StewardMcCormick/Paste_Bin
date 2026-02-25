package repository

import (
	"context"

	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type APIKeyRepository interface {
	Create(ctx context.Context, userId int64, key *domain.APIKey) (*domain.APIKey, error)
	RevokeKeyByUserId(ctx context.Context, userId int64) error
}

type TxUnitOfWork interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)

	UserRepository() UserRepository
	APIKeyRepository() APIKeyRepository
}

type NoTxUnitOfWork interface {
	UserRepository() UserRepository
	APIKeyRepository() APIKeyRepository
}

type pgxUnitOfWorkNoTx struct {
	pool *pgxpool.Pool
}

func (uw *pgxUnitOfWorkNoTx) UserRepository() UserRepository {
	return &userRepository{pool: uw.pool}
}

func (uw *pgxUnitOfWorkNoTx) APIKeyRepository() APIKeyRepository {
	return &apiKeyRepository{pool: uw.pool}
}

type pgxUnitOfWorkTX struct {
	tx pgx.Tx
}

func (uwt *pgxUnitOfWorkTX) Commit(ctx context.Context) error {
	return uwt.tx.Commit(ctx)
}

func (uwt *pgxUnitOfWorkTX) Rollback(ctx context.Context) {
	uwt.tx.Rollback(ctx)
}

func (uwt *pgxUnitOfWorkTX) UserRepository() UserRepository {
	return &userRepository{pool: uwt.tx}
}

func (uwt *pgxUnitOfWorkTX) APIKeyRepository() APIKeyRepository {
	return &apiKeyRepository{pool: uwt.tx}
}
