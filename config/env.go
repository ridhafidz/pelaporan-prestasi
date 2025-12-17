package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadENV() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system env")
	}
}

func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
