package user

import (
	"context"
	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	resultUser := &domain.User{}
	tx, err := r.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return nil, err
	}

	userQuery := `INSERT INTO users(username, password_hash, created_at) 
    				VALUES ($1, $2, $3) RETURNING *`
	err = tx.QueryRow(ctx, userQuery, user.Username, user.Password, user.CreatedAt).
		Scan(&resultUser.Id, &resultUser.Username, &resultUser.Password, &resultUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	apiKeyQuery := `INSERT INTO api_key(key_hash, user_id, created_at, expire_at, key_prefix) 
					VALUES ($1, $2, $3, $4, $5) RETURNING (expire_at)`
	err = tx.QueryRow(ctx, apiKeyQuery,
		user.APIKey.Key, resultUser.Id, user.CreatedAt, user.APIKey.ExpiresAt, user.APIKey.Prefix).
		Scan(&resultUser.APIKey.ExpiresAt)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	resultUser.APIKey.Key = user.APIKey.Key

	return resultUser, nil
}

func (r *repository) Exists(ctx context.Context, username string) (bool, error) {
	query := `SELECT * FROM users WHERE username = $1`
	row, err := r.pool.Query(ctx, query, username)
	defer row.Close()

	if err != nil {
		return false, err
	}

	if row.Next() {
		return true, nil
	}

	return false, nil
}

func (r *repository) GetAPIKeys(ctx context.Context, user *domain.User) ([]domain.APIKey, error) {
	return nil, nil
}
