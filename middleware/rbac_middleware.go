package middleware

import (
	"backend/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getClaims(c *fiber.Ctx) (*models.JWTClaims, error) {
	userToken := c.Locals("user")
	if userToken == nil {
		return nil, fiber.ErrUnauthorized
	}

	token := userToken.(*jwt.Token)
	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}

func OnlyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := getClaims(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		if claims.Role != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Admin only",
			})
		}

		return c.Next()
	}
}

func OnlyMahasiswa() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := getClaims(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		if claims.Role != "Mahasiswa" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Mahasiswa only",
			})
		}

		if claims.StudentID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Student ID not found in token",
			})
		}

		c.Locals("student_id", *claims.StudentID)

		return c.Next()
	}
}

func OnlyDosenWali() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := getClaims(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		if claims.Role != "DosenWali" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Dosen wali only",
			})
		}

		if claims.LecturerID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Lecturer ID not found in token",
			})
		}

		c.Locals("lecturer_id", *claims.LecturerID)

		return c.Next()
	}
}

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := getClaims(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		for _, p := range claims.Permissions {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Permission denied",
		})
	}
}
