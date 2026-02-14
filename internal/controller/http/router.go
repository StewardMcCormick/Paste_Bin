package http

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	handler := NewHandler()

	r.Use(LoggerMiddleware(logger))

	r.Get("/", handler.HelloHandler)

	return r
}
