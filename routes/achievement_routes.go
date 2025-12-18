package routes

import (
	"context"

	"backend/app/models"
	"backend/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// listAchievements godoc
// @Summary      List Achievements
// @Description  Get a list of achievement references for the logged-in student
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.AchievementReference
// @Router       /api/v1/achievements [get]
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

// getAchievementDetail godoc
// @Summary      Get Achievement Detail
// @Description  Retrieve full details of an achievement from MongoDB
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Mongo Achievement ID"
// @Security     ApiKeyAuth
// @Success      200  {object}  models.Achievement
// @Router       /api/v1/achievements/{id} [get]
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

// createAchievement godoc
// @Summary      Create Achievement
// @Description  Create a new achievement record (Mahasiswa only)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        request  body      models.Achievement  true  "Achievement Data"
// @Security     ApiKeyAuth
// @Success      201      {object}  map[string]string "Returns Mongo ID"
// @Router       /api/v1/achievements [post]
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

// updateAchievement godoc
// @Summary      Update Achievement
// @Description  Update draft achievement details (Mahasiswa only)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Mongo Achievement ID"
// @Param        request  body      models.Achievement true  "Updated Achievement Data"
// @Security     ApiKeyAuth
// @Success      200      {string}  string  "OK"
// @Router       /api/v1/achievements/{id} [put]
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

// deleteAchievement godoc
// @Summary      Delete Achievement
// @Description  Soft delete achievement from system (Mahasiswa only)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Mongo Achievement ID"
// @Security     ApiKeyAuth
// @Success      200  {string}  string  "OK"
// @Router       /api/v1/achievements/{id} [delete]
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

// submitAchievement godoc
// @Summary      Submit Achievement
// @Description  Change status from Draft to Submitted for verification
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Mongo Achievement ID"
// @Security     ApiKeyAuth
// @Success      200  {string}  string  "OK"
// @Router       /api/v1/achievements/{id}/submit [post]
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

// verifyAchievement godoc
// @Summary      Verify Achievement
// @Description  Approve achievement and assign points (Dosen Wali only)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Mongo Achievement ID"
// @Param        request  body      object  true  "Points data"
// @Security     ApiKeyAuth
// @Success      200      {string}  string  "OK"
// @Router       /api/v1/achievements/{id}/verify [post]
func verifyAchievement(
	refService service.AchievementReferenceService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals("user_id").(uuid.UUID)

		var body struct {
			Points float64 `json:"points"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid points format")
		}

		if err := refService.Verify(
			c.Context(),
			id,
			verifierID,
			body.Points,
		); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// rejectAchievement godoc
// @Summary      Reject Achievement
// @Description  Reject achievement with a note (Dosen Wali only)
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Mongo Achievement ID"
// @Param        request  body      object  true  "Rejection note"
// @Security     ApiKeyAuth
// @Success      200      {string}  string  "OK"
// @Router       /api/v1/achievements/{id}/reject [post]
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

// achievementHistory godoc
// @Summary      Get Achievement History
// @Description  Retrieve status history and metadata of an achievement
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Mongo Achievement ID"
// @Security     ApiKeyAuth
// @Success      200  {object}  models.AchievementReference
// @Router       /api/v1/achievements/{id}/history [get]
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

// addAttachment godoc
// @Summary      Add Attachment
// @Description  Upload/Add file URL to an achievement
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Mongo Achievement ID"
// @Param        request  body      models.Attachment  true  "Attachment data"
// @Security     ApiKeyAuth
// @Success      201      {string}  string             "Created"
// @Router       /api/v1/achievements/{id}/attachments [post]
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
