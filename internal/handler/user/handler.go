package user

import (
	"context"

	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
)

type UseCase interface {
	Registration(ctx context.Context, user *dto.UserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, user *dto.UserRequest) (*dto.APIKeyResponse, error)
}

type httpHandlers struct {
	authUseCase UseCase
}

func NewHandler(authUseCase UseCase) *httpHandlers {
	return &httpHandlers{
		authUseCase: authUseCase,
	}
}
