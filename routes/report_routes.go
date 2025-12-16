package routes

import (
	"net/http"
	"strconv"

	"backend/app/models"
	"backend/app/services"
	"backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterReportRoutes registers reporting endpoints
func RegisterReportRoutes(r *gin.Engine, achievementService service.AchievementService, achRefService service.AchievementReferenceService, studentService service.StudentService) {
	rep := r.Group("/reports")
	rep.Use(middleware.JWTMiddleware())

	// Student report: student info + achievements + achievement references
	rep.GET("/student/:id", middleware.RBAC("admin", "dosen", "mahasiswa"), func(c *gin.Context) {
		id := c.Param("id")

		// student details (studentService expects string id)
		student, err := studentService.FindByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "student not found"})
			return
		}

		// achievements (stored in Mongo, service accepts string studentID)
		achievements, err := achievementService.GetByStudentID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		// achievement references (postgres, repo expects uuid)
		var refs []models.AchievementReference
		if uid, err := uuid.Parse(id); err == nil {
			refs, _ = achRefService.GetByStudentID(c, uid)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"student": student, "achievements": achievements, "references": refs}})
	})

	// Simple summary: count achievements for a year (optional)
	rep.GET("/summary", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		yearStr := c.Query("year")
		if yearStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "year query param required"})
			return
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid year"})
			return
		}

		// naive summary: fetch all achievements for that year by scanning student IDs won't scale
		// Here we call FindByTag("year:<year>") if tag convention used; fallback to empty
		tag := "year:" + strconv.Itoa(year)
		results, err := achievementService.FindByTag(c, tag, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"year": year, "count": len(results)}})
	})
}
