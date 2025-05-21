package main

import (
	"exceapp/internals/router"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	// "github.com/joho/godotenv"
)

func main() {

	route := chi.NewRouter()
	route.Mount("/api/auth", router.Auth())
	fmt.Println(os.Getenv("DB_USER"))
	fmt.Println("docker runnning....🎉")
	log.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", route)
}
