package handler

import (
	"context"
	"encoding/json"
	"github.com/StewardMcCormick/Paste_Bin/internal/util"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendError(ctx context.Context, w http.ResponseWriter, status int, message error) {
	log := util.GetLoggerFromCtx(ctx)

	response := ErrorResponse{status, message.Error()}
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Error(message.Error())
}
