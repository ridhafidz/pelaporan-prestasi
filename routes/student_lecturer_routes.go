package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/service"
)

func StudentList(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.GetStudents(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
}

func StudentGetByID(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		data, err := svc.GetStudentDetail(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return c.JSON(data)
	}
}

func StudentAchievements(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		data, err := svc.GetStudentAchievements(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return c.JSON(data)
	}
}

func StudentUpdateAdvisor(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		var req models.UpdateAdvisorRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		if err := svc.UpdateAdvisor(c.Context(), id, req.AdvisorID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(fiber.Map{
			"message": "advisor updated successfully",
		})
	}
}

func LecturerList(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.GetLecturers(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
}

func LecturerAdvisees(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		data, err := svc.GetLecturerAdvisees(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(data)
	}
}
