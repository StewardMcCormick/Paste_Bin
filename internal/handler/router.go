package handler

import (
	"github.com/StewardMcCormick/Paste_Bin/config/cfg_util"
	midd "github.com/StewardMcCormick/Paste_Bin/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)
}

func NewRouter(userHandler UserHandler, logger *zap.Logger, env cfgUtil.Env) http.Handler {
	r := chi.NewRouter()

	r.Use(midd.LoggerMiddleware(logger))
	r.Use(midd.RecovererMiddleware)
	r.Use(midd.EnvironmentalMiddleware(env))

	r.Post("/user", userHandler.Registration)

	return r
}
