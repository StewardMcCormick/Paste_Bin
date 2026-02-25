package repository

import (
	"context"

	"github.com/StewardMcCormick/Paste_Bin/internal/adapter/postgres"
	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
)

type userRepository struct {
	pool postgres.DBTX
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	log := appctx.GetLogger(ctx)
	resultUser := &domain.User{}

	userQuery := `INSERT INTO users(username, password_hash, created_at) 
    				VALUES ($1, $2, $3) RETURNING *`
	err := r.pool.QueryRow(ctx, userQuery, user.Username, user.Password, user.CreatedAt).
		Scan(&resultUser.Id, &resultUser.Username, &resultUser.Password, &resultUser.CreatedAt)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resultUser, nil
}

func (r *userRepository) Exists(ctx context.Context, username string) (bool, error) {
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
