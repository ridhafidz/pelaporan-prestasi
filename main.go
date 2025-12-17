package main

import (
	"log"
	// "backend/config" 
	"backend/database" 

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	log.Println("⏳ Mencoba menghubungkan database...")
	database.ConnectPostgres()
	database.ConnectMongo()

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