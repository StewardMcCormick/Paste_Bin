package handler

import (
	"net/http"

	mid "github.com/StewardMcCormick/Paste_Bin/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	NotFound(w http.ResponseWriter, r *http.Request)
	MethodNotAllowed(w http.ResponseWriter, r *http.Request)
	Registration(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Hello(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	userHandler UserHandler,
	logMid mid.Logging,
	recovererMid mid.Recoverer,
	envMid mid.Environmental,
	validMid mid.JSONValidation,
	authMid mid.Auth,
) http.Handler {
	r := chi.NewRouter()

	r.Use(logMid.Handler)
	r.Use(recovererMid.Handler)
	r.Use(envMid.Handler)
	r.Use(validMid.Handler)

	r.NotFound(userHandler.NotFound)
	r.MethodNotAllowed(userHandler.MethodNotAllowed)

	r.Post("/registration", userHandler.Registration)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMid.Handler)
		r.Get("/hello", userHandler.Hello)
	})

	return r
}
