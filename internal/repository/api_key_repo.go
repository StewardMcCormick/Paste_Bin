package repository

import (
	"context"

	"github.com/StewardMcCormick/Paste_Bin/internal/adapter/postgres"
	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
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

func (r *apiKeyRepository) Exists(ctx context.Context, keyHash string) (bool, error) {
	return false, nil
}
