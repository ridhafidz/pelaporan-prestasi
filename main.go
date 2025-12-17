package main

import (
	"log"
	"os"

	"backend/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	app := config.InitApp()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Server running on port:", port)
	log.Fatal(app.App.Listen(":" + port))
}
