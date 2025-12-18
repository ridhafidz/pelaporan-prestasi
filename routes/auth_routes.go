package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/service"
	"backend/middleware"
)

// processLogin godoc
// @Summary      User Login
// @Description  Authenticate user and return access & refresh tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login Credentials"
// @Success      200      {object}  models.LoginResponse
// @Failure      401      {object}  map[string]string "Invalid email or password"
// @Router       /api/v1/auth/login [post]
func processLogin(s service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.LoginRequest)

		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		authResponse, err := s.Login(c.Context(), *req)
		if err != nil {
			if err.Error() == "invalid credentials" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": authResponse})
	}
}

// processRefreshToken godoc
// @Summary      Refresh Token
// @Description  Get a new access token using a valid refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200      {object}  models.RefreshTokenResponse
// @Failure      401      {object}  map[string]string
// @Router       /api/v1/auth/refresh [post]
func processRefreshToken(s service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Authorization header missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization format",
			})
		}

		refreshToken := parts[1]

		resp, err := s.RefreshToken(c.Context(), refreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// processLogout godoc
// @Summary      User Logout
// @Description  Invalidate the current refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200      {object}  map[string]string "Logged out successfully"
// @Router       /api/v1/auth/logout [post]
func processLogout(s service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization format",
			})
		}

		req := models.LogoutRequest{
			RefreshToken: parts[1],
		}

		if err := s.Logout(c.Context(), req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Logged out successfully",
		})
	}
}

// processGetProfile godoc
// @Summary      Get Current Profile
// @Description  Retrieve the logged-in user's profile information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200      {object}  models.UserData
// @Failure      404      {object}  map[string]string "User not found"
// @Router       /api/v1/auth/profile [get]
func processGetProfile(s service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userIDVal := c.Locals(middleware.UserIDKey)
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
