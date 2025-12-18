package main

import (
	"log"
	"os"

	"backend/config"

	"github.com/joho/godotenv"
)

// @title           Sistem Pelaporan Prestasi API
// @version         1.0
// @description     API Server untuk manajemen dan pelaporan prestasi mahasiswa.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
