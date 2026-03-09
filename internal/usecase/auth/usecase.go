package auth

import (
	"context"
	"time"

	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	"github.com/StewardMcCormick/Paste_Bin/internal/repository"
)

type Config struct {
	APIKeyExpireDuration time.Duration `yaml:"api_key_expires_time" env-default:"168h"`
}

type UnitOfWorkFactory interface {
	Exec(ctx context.Context) repository.NoTxUnitOfWork
	Begin(ctx context.Context) (repository.TxUnitOfWork, error)
}

type Security interface {
	HashPassword(password string) (string, error)
	HashAPIKey(key string) string
	GenerateAPIKey(ctx context.Context) (keyPrefix string, key string, err error)
	CompareHashAndPassword(hash string, pass string) bool
}

type Validator interface {
	Validate(request *dto.UserRequest) error
}

type GeneratedAPIKey struct {
	Prefix string
	Key    string
	Hash   string
}

type UseCase struct {
	uow          UnitOfWorkFactory
	securityUtil Security
	valid        Validator
	cfg          Config
}

func NewUseCase(uow UnitOfWorkFactory, securityUtil Security, valid Validator, cfg Config) *UseCase {
	return &UseCase{
		uow:          uow,
		securityUtil: securityUtil,
		valid:        valid,
		cfg:          cfg,
	}
}

func (uc *UseCase) generateNewKey(ctx context.Context) (GeneratedAPIKey, error) {
	prefix, apiKey, err := uc.securityUtil.GenerateAPIKey(ctx)
	if err != nil {
		return GeneratedAPIKey{}, err
	}
	hashedKey := uc.securityUtil.HashAPIKey(apiKey)

	return GeneratedAPIKey{Prefix: prefix, Key: apiKey, Hash: hashedKey}, nil
}
