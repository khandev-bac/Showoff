package main

import (
	"exceapp/cmd/config"
	"exceapp/internals/handler"
	"exceapp/internals/model"
	"exceapp/internals/repo"
	"exceapp/internals/router"
	"exceapp/internals/service"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not loaded properly")
	}
	config.ConnectDB()
	db := config.DB
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("❌ AutoMigrate failed: " + err.Error())
	}
	UserRepo := *repo.NewUserRepo(db)
	UserService := *service.NewUserService(&UserRepo)
	UserHandler := *handler.NewUserHandler(&UserService)
	r := chi.NewRouter()
	r.Mount("/api/auth", router.AuthRoutes(&UserHandler))
	// url := google.GetLoginUrl("khan")
	// fmt.Println(url)
	fmt.Println("docker runnning....🎉")
	log.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
