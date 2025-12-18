package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/service"
)

var validate = validator.New()

// processCreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user record in the system (Admin only)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateUserRequest  true  "User Data"
// @Success      201      {object}  models.UserResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string "Email/Username already taken"
// @Security     ApiKeyAuth
// @Router       /api/v1/users [post]
func processCreateUser(s service.UserService) fiber.Handler {
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

// processGetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve detailed information of a specific user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  models.UserResponse
// @Failure      404  {object}  map[string]string "User not found"
// @Security     ApiKeyAuth
// @Router       /api/v1/users/{id} [get]
func processGetUserByID(s service.UserService) fiber.Handler {
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

// processUpdateUser godoc
// @Summary      Update user
// @Description  Update user personal information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "User UUID"
// @Param        request  body      models.UpdateUserRequest  true  "Updated User Data"
// @Success      200      {object}  map[string]string         "User updated successfully"
// @Security     ApiKeyAuth
// @Router       /api/v1/users/{id} [put]
func processUpdateUser(s service.UserService) fiber.Handler {
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

// processDeleteUser godoc
// @Summary      Delete user
// @Description  Remove user from the system
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string             true  "User UUID"
// @Success      200  {object}  map[string]string  "User deleted successfully"
// @Security     ApiKeyAuth
// @Router       /api/v1/users/{id} [delete]
func processDeleteUser(s service.UserService) fiber.Handler {
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

// processGetAllUsers godoc
// @Summary      List all users
// @Description  Retrieve a paginated list of all registered users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number (default: 1)"
// @Param        limit  query     int  false  "Items per page (default: 10)"
// @Success      200    {object}  map[string]interface{}
// @Security     ApiKeyAuth
// @Router       /api/v1/users [get]
func processGetAllUsers(s service.UserService) fiber.Handler {
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

// processUpdateUserRole godoc
// @Summary      Update user role
// @Description  Change the access role of a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "User UUID"
// @Param        request  body      models.UpdateUserRoleRequest  true  "Role Data"
// @Success      200      {object}  map[string]string             "User role updated successfully"
// @Security     ApiKeyAuth
// @Router       /api/v1/users/{id}/role [put]
func processUpdateUserRole(s service.UserService) fiber.Handler {
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
