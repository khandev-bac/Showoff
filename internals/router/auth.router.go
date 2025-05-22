package router

import (
	"exceapp/cmd/config"
	"exceapp/internals/handler"
	"exceapp/internals/middleware"
	"exceapp/internals/repo"
	"exceapp/internals/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func Auth() http.Handler {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not loaded properly")
	}
	r := chi.NewRouter()
	config.InitDB()
	repo := repo.NewUserRepo(config.DB)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)
	r.Post("/signup", handler.Signup)
	r.Post("/login", handler.Login)
	r.Get("/ok", handler.Check)
	r.Get("/google-login", handler.GoogleLogin)
	r.Get("/google-callback", handler.GoogleCallback)
	r.Get("/refreshToken", handler.RefreshTokenHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.Auth)
		protected.Get("/logout", handler.Logout)
		protected.Get("/user", handler.GetUserInfo)
		protected.Post("/upload-profile", handler.UploadProfilePic)
	})
	return r
}
