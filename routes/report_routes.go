package routes

import (
	"backend/app/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// handleGetStatistics godoc
// @Summary      Get Achievement Statistics
// @Description  Retrieve statistics based on role (Admin: Global, Dosen Wali: Advisees, Mahasiswa: Own).
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start  query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end    query     string  false  "End date (YYYY-MM-DD)"
// @Security     ApiKeyAuth
// @Success      200    {object}  models.AchievementStatisticsResponse
// @Failure      403    {object}  map[string]string "Invalid role or access denied"
// @Router       /api/v1/reports/statistics [get]
func handleGetStatistics(s service.ReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string) 
		userIDStr := c.Locals("userId").(string)
		userID, _ := uuid.Parse(userIDStr)

		startStr := c.Query("start")
		endStr := c.Query("end")

		start, _ := time.Parse("2006-01-02", startStr)
		end, _ := time.Parse("2006-01-02", endStr)
		
		if end.IsZero() {
			end = time.Now()
		}

		stats, err := s.GetStatistics(c.Context(), userRole, userID, start, end)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   stats,
		})
	}
}

// handleGetStudentReport godoc
// @Summary      Get Student Individual Report
// @Description  Get total points and student details (Accessible by Mhs for own, Dosen for advisees).
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Student UUID"
// @Security     ApiKeyAuth
// @Success      200  {object}  models.StudentReportResponse
// @Failure      403  {object}  map[string]string "Access denied"
// @Router       /api/v1/reports/student/{id} [get]
func handleGetStudentReport(s service.ReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		userIDStr := c.Locals("userId").(string)
		requesterID, _ := uuid.Parse(userIDStr)
		targetID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid student id"})
		}

		report, err := s.GetStudentReport(c.Context(), userRole, requesterID, targetID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   report,
		})
	}
}