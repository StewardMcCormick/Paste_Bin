package user

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/usecase/auth"
)

type httpHandlers struct {
	authUseCase *auth.UseCase
}

func NewHandler(authUseCase *auth.UseCase) *httpHandlers {
	return &httpHandlers{
		authUseCase: authUseCase,
	}
}
