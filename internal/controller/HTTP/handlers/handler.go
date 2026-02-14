package handlers

import "net/http"

type Handlers interface {
	HelloHandler(w http.ResponseWriter, r *http.Request)
}

type httpHandlers struct {
}

func NewHandler() Handlers {
	return &httpHandlers{}
}
