package routes

import (
	"github.com/gofiber/fiber/v2"

	models "backend/app/models/postgree"
	services "backend/app/services/postgree"

	"github.com/google/uuid"
)

func processLogin(s services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.LoginRequest)

		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		authResponse, err := s.Login(c.Context(), req)
		if err != nil {
			if err.Error() == "invalid credentials" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": authResponse})
	}
}

func processRefreshToken(s services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.RefreshTokenRequest)

		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		resp, err := s.RefreshToken(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func processLogout(s services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing token"})
		}

		if err := s.Logout(c.Context(), authHeader); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Logged out successfully"})
	}
}

func processGetProfile(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userIDVal := c.Locals("userID")
		if userIDVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			uidStr, okStr := userIDVal.(string)
			if okStr {
				parsed, err := uuid.Parse(uidStr)
				if err == nil {
					userID = parsed
				} else {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid User ID format in token"})
				}
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid User ID in context"})
			}
		}

		user, err := s.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
	}
}
