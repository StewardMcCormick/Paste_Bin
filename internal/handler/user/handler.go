package user

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/usecase/user"
)

type httpHandlers struct {
	UserUseCase *user.UseCase
}

func NewHandler(userUseCase *user.UseCase) *httpHandlers {
	return &httpHandlers{
		UserUseCase: userUseCase,
	}
}
