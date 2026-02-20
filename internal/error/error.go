package error

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StewardMcCormick/Paste_Bin/internal/util"
	"net/http"
)

var (
	InternalError = errors.New("internal error")

	// Domain error
	UserAlreadyExists = errors.New("user already exists")
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendHTTPError(ctx context.Context, w http.ResponseWriter, status int, message error) {
	log := util.GetLoggerFromCtx(ctx)

	response := Response{status, message.Error()}
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info(message.Error())
}
