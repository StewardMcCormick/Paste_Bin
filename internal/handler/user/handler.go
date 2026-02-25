package user

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/usecase/auth"
)

type httpHandlers struct {
	UserUseCase *auth.UseCase
}

func NewHandler(userUseCase *auth.UseCase) *httpHandlers {
	return &httpHandlers{
		UserUseCase: userUseCase,
	}
}
