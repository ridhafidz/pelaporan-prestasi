package main

import (
	"log"
	// "backend/config" 
	"backend/database" 

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	// 2. Test Koneksi Database
	log.Println("⏳ Mencoba menghubungkan database...")
	database.ConnectPostgres()
	database.ConnectMongo()

	// 3. Setup Fiber (Hanya untuk menjaga aplikasi tetap jalan)
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "success",
			"message":  "Database connections are healthy!",
			"postgres": "connected",
			"mongo":    "connected",
		})
	})

	// Jalankan server
	log.Fatal(app.Listen(":3000"))
}