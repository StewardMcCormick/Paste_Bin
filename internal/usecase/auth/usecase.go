package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	"github.com/StewardMcCormick/Paste_Bin/internal/repository"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	"go.uber.org/zap"
)

type Config struct {
	APIKeyExpireDuration time.Duration `yaml:"api_key_expires_time" env-default:"168h"`
	MaxUserKeysNum       int           `yaml:"max_user_keys_num" env-default:"5"`
}

type UnitOfWorkFactory interface {
	Exec(ctx context.Context) (repository.NoTxUnitOfWork, error)
	Begin(ctx context.Context) (repository.TxUnitOfWork, error)
}

type SecurityUtil interface {
	HashPassword(password string) (string, error)
	HashAPIKey(key string) string
	GenerateAPIKey(ctx context.Context) (keyPrefix string, key string, err error)
	CompareHashAndPassword(hash string, pass string) bool
}

type Validator interface {
	Validate(request *dto.UserRequest) error
}

type UseCase struct {
	uow          UnitOfWorkFactory
	securityUtil SecurityUtil
	valid        Validator
	cfg          Config
}

func NewUseCase(uow UnitOfWorkFactory, securityUtil SecurityUtil, valid Validator, cfg Config) *UseCase {
	return &UseCase{
		uow:          uow,
		securityUtil: securityUtil,
		valid:        valid,
		cfg:          cfg,
	}
}

func (uc *UseCase) Registration(ctx context.Context, user *dto.UserRequest) (*dto.UserResponse, error) {
	log := appctx.GetLogger(ctx)
	if err := uc.valid.Validate(user); err != nil {
		return nil, err
	}

	tx, err := uc.uow.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w - database error", err)
	}
	defer tx.Rollback(ctx)

	exists, err := tx.UserRepository().Exists(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("%w - database error", errs.InternalError)
	}
	if exists {
		return nil, errs.UserAlreadyExists
	}

	hashedPass, err := uc.securityUtil.HashPassword(user.Password)
	if err != nil {
		log.Warn(fmt.Sprintf("%s - password hashing error", err.Error()))
		return nil, fmt.Errorf("%w - password hashing error", errs.InternalError)
	}
	prefix, apiKey, err := uc.securityUtil.GenerateAPIKey(ctx)
	if err != nil {
		log.Warn(fmt.Sprintf("%s - API key generation error", err.Error()))
		return nil, fmt.Errorf("%w - API key generation error", errs.InternalError)
	}
	hashedKey := uc.securityUtil.HashAPIKey(apiKey)

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
			Key:       hashedKey,
			Prefix:    prefix,
			CreatedAt: now,
			ExpiresAt: now.Add(uc.cfg.APIKeyExpireDuration),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w - API key create error", errs.InternalError)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w - tx commit error", errs.InternalError)
	}

	createdApiKey.Key = apiKey
	createdUser.APIKey = *createdApiKey

	log.Info(
		"user registered",
		zap.String("username", createdUser.Username),
	)
	return createdUser.ToResponse(), nil
}

//func (uc *UseCase) Login(ctx context.Context, auth *dto.UserRequest) (*dto.UserResponse, error) {
//	//exists, err := uc.repo.Exists(ctx, auth.Username)
//	//if err != nil {
//	//	return nil, fmt.Errorf("%w - database error", errs.InternalError)
//	//}
//	//if !exists {
//	//	return nil, errs.UserNotFound
//	//}
//	//
//	//userFromDb, err := uc.repo.GetUserWithAvailableKeyByUsername(ctx, auth.Username)
//	//if userFromDb == nil {
//	//	if err != nil {
//	//		return nil, fmt.Errorf("%w - database error", errs.InternalError)
//	//	}
//	//	return nil, errs.UserNotFound
//	//}
//	//
//	//if !uc.securityUtil.CompareHashAndPassword(userFromDb.Password, auth.Password) {
//	//	return nil, errs.Unforbidden
//	//}
//	//
//	//return userFromDb.ToResponse(), nil
//}
