package router

import (
	"exceapp/internals/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(handler *handler.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/google-login", handler.GoogleLogin)
	r.Get("/google-callback", handler.GoogleCallback)
	r.Post("/login", handler.Login)
	r.Post("/signup", handler.Signup)
	r.Get("/ok", handler.Check)
	return r
}
