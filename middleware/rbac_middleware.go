package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func OnlyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user")
		if userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Unauthorized access",
			})
		}

		token := userToken.(*jwt.Token)
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. Admins only.",
			})
		}

		return c.Next()
	}
}

func OnlyMahasiswa() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user")
		if userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Unauthorized access",
			})
		}

		token := userToken.(*jwt.Token)
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "Mahasiswa" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. Mahasiswa only.",
			})
		}

		return c.Next()
	}
}

func OnlyDosenWali() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user")
		if userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Unauthorized access",
			})
		}

		token := userToken.(*jwt.Token)
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || claims["role"] != "DosenWali" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. Dosen wali only.",
			})
		}

		return c.Next()
	}
}
