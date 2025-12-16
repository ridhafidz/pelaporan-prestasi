package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	models "backend/app/models/postgree"
	services "backend/app/services/postgree"
)

var validate = validator.New()

func processCreateUser(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(models.CreateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}
		
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		createdUser, err := s.CreateUser(c.Context(), req)
		if err != nil {
			if err.Error() == "email already registered" || err.Error() == "username already taken" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": createdUser})
	}
}

func processGetUserByID(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		userID, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
		}

		user, err := s.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
	}
}

func processUpdateUser(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		userID, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
		}

		req := new(models.UpdateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}
		
		if err := s.UpdateUser(c.Context(), userID, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User updated successfully"})
	}
}

func processDeleteUser(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		userID, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
		}

		if err := s.DeleteUser(c.Context(), userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User deleted successfully"})
	}
}

func processGetAllUsers(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		users, err := s.GetAllUsers(c.Context(), page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   users,
			"meta":   fiber.Map{"page": page, "limit": limit},
		})
	}
}

func processUpdateUserRole(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		userID, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID format"})
		}

		req := new(models.UpdateUserRoleRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		if err := s.UpdateUserRole(c.Context(), userID, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User role updated successfully"})
	}
}