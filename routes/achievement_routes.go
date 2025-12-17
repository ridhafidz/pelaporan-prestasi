package routes

import (
	"context"

	"backend/app/models"
	"backend/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func listAchievements(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Locals("student_id").(uuid.UUID)

		data, err := refService.GetByStudentID(
			context.Background(),
			studentID,
			10,
			0,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(data)
	}
}

func getAchievementDetail(
	achievementService service.AchievementService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		data, err := achievementService.GetByID(context.Background(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return c.JSON(data)
	}
}

func createAchievement(
	achievementService service.AchievementService,
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload models.Achievement

		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		studentID := c.Locals("student_id").(uuid.UUID)
		payload.StudentID = studentID.String()

		mongoID, err := achievementService.Create(
			context.Background(),
			&payload,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err := refService.Create(
			context.Background(),
			studentID,
			mongoID,
		); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": mongoID,
		})
	}
}

func updateAchievement(
	achievementService service.AchievementService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var payload models.Achievement
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := achievementService.Update(
			context.Background(),
			id,
			&payload,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func deleteAchievement(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := refService.Delete(
			context.Background(),
			id,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func submitAchievement(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := refService.Submit(context.Background(), id); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func verifyAchievement(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals("user_id").(uuid.UUID)

		if err := refService.Verify(
			context.Background(),
			id,
			verifierID,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func rejectAchievement(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var body struct {
			Note string `json:"note"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := refService.Reject(
			context.Background(),
			id,
			body.Note,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func achievementHistory(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		data, err := refService.GetByMongoID(context.Background(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return c.JSON(data)
	}
}

func addAttachment(
	achievementService service.AchievementService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var attachment models.Attachment
		if err := c.BodyParser(&attachment); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := achievementService.AddAttachment(
			context.Background(),
			id,
			attachment,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusCreated)
	}
}
