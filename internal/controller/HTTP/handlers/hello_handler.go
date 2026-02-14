package handlers

import (
	"github.com/StewardMcCormick/Paste_Bin/pkg/render"
	"net/http"
)

func (h *httpHandlers) HelloHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, struct {
		Message string `json:"message"`
	}{"Hello world!"})
}
