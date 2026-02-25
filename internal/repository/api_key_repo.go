package repository

import (
	"context"

	"github.com/StewardMcCormick/Paste_Bin/internal/adapter/postgres"
	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
)

type apiKeyRepository struct {
	pool postgres.DBTX
}

func (r *apiKeyRepository) Create(ctx context.Context, userId int64, key *domain.APIKey) (*domain.APIKey, error) {
	log := appctx.GetLogger(ctx)

	query := `INSERT INTO api_key(key_hash, user_id, created_at, expire_at, key_prefix) 
					VALUES ($1, $2, $3, $4, $5) RETURNING (expire_at)`
	_, err := r.pool.Exec(ctx, query, key.Key, userId, key.CreatedAt, key.ExpiresAt, key.Prefix)

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	return key, nil
}

func (r *apiKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	log := appctx.GetLogger(ctx)

	query := `SELECT (key_hash, created_at, expire_at, key_prefix) FROM api_key WHERE key_hash=$1`
	rows, err := r.pool.Query(ctx, query, keyHash)
	defer rows.Close()

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if rows.Next() {
		resultKey := &domain.APIKey{}
		if err = rows.Scan(&resultKey.Key, &resultKey.CreatedAt, &resultKey.ExpiresAt, &resultKey.Prefix); err != nil {
			return nil, err
		}

		return resultKey, nil
	}

	return nil, errs.APIKeyNotFound
}

func (r *apiKeyRepository) RevokeKeyByUserId(ctx context.Context, userId int64) error {
	log := appctx.GetLogger(ctx)

	query := `DELETE FROM api_key WHERE user_id=$1`
	_, err := r.pool.Exec(ctx, query, userId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
