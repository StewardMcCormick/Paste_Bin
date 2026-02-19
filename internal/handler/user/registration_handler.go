package user

import (
	"encoding/json"
	"errors"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	"github.com/StewardMcCormick/Paste_Bin/internal/handler"
	"net/http"
)

func (h *httpHandlers) Registration(w http.ResponseWriter, r *http.Request) {
	var userRequest dto.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&userRequest)

	user, err := h.UserUseCase.Registration(r.Context(), &userRequest)
	if err != nil {
		if errors.Is(err, errs.UserAlreadyExists) {
			handler.SendError(r.Context(), w, http.StatusConflict, err)
			return
		} else if errors.Is(err, errs.InternalError) {
			handler.SendError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}
		handler.SendError(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		handler.SendError(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
}
