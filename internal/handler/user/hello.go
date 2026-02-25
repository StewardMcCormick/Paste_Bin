package user

import (
	"net/http"

	"github.com/StewardMcCormick/Paste_Bin/pkg/render"
)

func (h *httpHandlers) Hello(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, struct {
		Message string `json:"message"`
	}{"hello, world!"})
}
