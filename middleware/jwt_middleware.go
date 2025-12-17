package middleware

import (
	"strings"

	"backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey      = "userID"
	UsernameKey    = "username"
	RoleKey        = "role"
	PermissionsKey = "permissions"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "missing or invalid Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid authorization format",
			})
		}

		tokenString := parts[1]

		claims, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid or expired token",
			})
		}

		c.Locals(UserIDKey, claims.UserID)
		c.Locals(UsernameKey, claims.Username)
		c.Locals(RoleKey, claims.Role)
		c.Locals(PermissionsKey, claims.Permissions)

		return c.Next()
	}
}
