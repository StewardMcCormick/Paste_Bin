package handlers

import "net/http"

type Handler interface {
	HelloHandler(w http.ResponseWriter, r *http.Request)
}

type httpHandler struct {
}

func NewHandler() Handler {
	return &httpHandler{}
}
