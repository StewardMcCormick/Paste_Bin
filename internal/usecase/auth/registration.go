package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	"go.uber.org/zap"
)

func (uc *UseCase) Registration(ctx context.Context, user *dto.UserRequest) (*dto.UserResponse, error) {
	log := appctx.GetLogger(ctx)
	if err := uc.valid.Validate(user); err != nil {
		log.Debug(fmt.Sprintf("validation error - %v", err))
		return nil, err
	}

	tx, err := uc.uow.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w - tx beggining error", errs.InternalError)
	}
	defer tx.Rollback(ctx)

	userFromDb, err := tx.UserRepository().GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("%w - database error", errs.InternalError)
	}
	if userFromDb != nil {
		return nil, errs.UserAlreadyExists
	}

	hashedPass, err := uc.securityUtil.HashPassword(user.Password)
	if err != nil {
		log.Warn(fmt.Sprintf("%s - password hashing error", err.Error()))
		return nil, fmt.Errorf("%w - password hashing error", errs.InternalError)
	}

	key, err := uc.generateNewKey(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("%s - API key generation error", err.Error()))
		return nil, fmt.Errorf("%w - API key generation error", errs.InternalError)
	}

	now := time.Now()
	createdUser, err := tx.UserRepository().Create(
		ctx, &domain.User{
			Username:  user.Username,
			Password:  hashedPass,
			CreatedAt: now,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w - user create error", errs.InternalError)
	}

	createdApiKey, err := tx.APIKeyRepository().Create(
		ctx, createdUser.Id, &domain.APIKey{
			UserId:    createdUser.Id,
			Key:       key.Hash,
			Prefix:    key.Prefix,
			CreatedAt: now,
			ExpiresAt: now.Add(uc.cfg.APIKeyExpireDuration),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w - API key create error", errs.InternalError)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("%w - tx commit error", errs.InternalError)
	}

	createdApiKey.Key = key.Key
	createdUser.APIKey = *createdApiKey

	log.Info(
		"user registered",
		zap.String("username", createdUser.Username),
	)
	return createdUser.ToResponse(), nil
}
