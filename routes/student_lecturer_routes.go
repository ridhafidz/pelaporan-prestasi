package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/service"
)

// StudentList godoc
// @Summary      List All Students
// @Description  Get a list of all students with their advisors (Admin only)
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.StudentDetailResponse
// @Router       /api/v1/students [get]
func StudentList(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.GetStudents(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
}

// StudentGetByID godoc
// @Summary      Get Student Detail
// @Description  Get detailed information of a student by their UUID
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Student UUID"
// @Security     ApiKeyAuth
// @Success      200  {object}  models.StudentDetailResponse
// @Router       /api/v1/students/{id} [get]
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

// StudentAchievements godoc
// @Summary      Get Student Achievements
// @Description  List all achievement status and points for a specific student
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Student UUID"
// @Security     ApiKeyAuth
// @Success      200  {array}   models.StudentAchievementResponse
// @Router       /api/v1/students/{id}/achievements [get]
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

// StudentUpdateAdvisor godoc
// @Summary      Assign/Update Advisor
// @Description  Assign a lecturer as an advisor to a student (Admin only)
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "Student UUID"
// @Param        request  body      models.UpdateAdvisorRequest  true  "Lecturer ID"
// @Security     ApiKeyAuth
// @Success      200      {object}  map[string]string "message: advisor updated successfully"
// @Router       /api/v1/students/{id}/advisor [put]
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

// LecturerList godoc
// @Summary      List All Lecturers
// @Description  Get a list of all lecturers in the system
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.LecturerDetailResponse
// @Router       /api/v1/lecturers [get]
func LecturerList(svc service.StudentLecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.GetLecturers(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
}

// LecturerAdvisees godoc
// @Summary      Get Lecturer Advisees
// @Description  List all students assigned to a specific lecturer
// @Tags         Students & Lecturers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lecturer UUID"
// @Security     ApiKeyAuth
// @Success      200  {array}   models.StudentDetailResponse
// @Router       /api/v1/lecturers/{id}/advisees [get]
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
