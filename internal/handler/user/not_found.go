package user

import (
	"net/http"

	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
)

func (h *httpHandlers) NotFound(w http.ResponseWriter, r *http.Request) {
	errs.SendAppError(r.Context(), w, http.StatusNotFound, errs.PageNotFound)
}
