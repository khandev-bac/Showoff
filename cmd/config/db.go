package config

import (
	"exceapp/internals/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found!")
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required env var: %s", key)
	}
	return val
}

func InitDB() {
	LoadEnv()

	host := getEnv("DB_HOST")
	port := getEnv("DB_PORT")
	user := getEnv("DB_USER")
	password := getEnv("DB_PASSWORD")
	dbname := getEnv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	log.Println("Connecting to DB with:")
	log.Printf("Host: %s | DB: %s | User: %s\n", host, dbname, user)

	var db *gorm.DB
	var err error

	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("[Attempt %d] Failed to connect to DB: %v\n", i, err)
		time.Sleep(time.Duration(i*2) * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to DB after %d attempts: %v", maxRetries, err)
	}

	DB = db
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("migration failed")
	}
	log.Println("✅ Successfully connected to the database")
}
