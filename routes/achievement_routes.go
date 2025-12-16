package routes

import (
	"net/http"

	"backend/app/models"
	"backend/app/services"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAchievementRoutes registers achievement-related HTTP routes
func RegisterAchievementRoutes(r *gin.Engine, achievementService service.AchievementService) {
	a := r.Group("/achievements")
	a.Use(middleware.JWTMiddleware())

	// Create achievement (admin/dosen)
	a.POST("/", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		var req models.Achievement
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		created, err := achievementService.Create(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "success", "data": created})
	})

	// List achievements by student_id or tag
	a.GET("/", middleware.RBAC("admin", "dosen", "mahasiswa"), func(c *gin.Context) {
		studentID := c.Query("student_id")
		tag := c.Query("tag")

		if studentID != "" {
			results, err := achievementService.GetByStudentID(c, studentID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
			return
		}

		if tag != "" {
			// optional limit param
			limit := int64(10)
			if l := c.Query("limit"); l != "" {
				// ignore parse errors and fall back to default
			}

			results, err := achievementService.FindByTag(c, tag, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "provide student_id or tag as query parameter"})
	})

	// Get achievement by id
	a.GET(":id", middleware.RBAC("admin", "dosen", "mahasiswa"), func(c *gin.Context) {
		id := c.Param("id")
		res, err := achievementService.GetByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res})
	})

	// Update achievement (admin/dosen)
	a.PUT(":id", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		id := c.Param("id")

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if err := achievementService.Update(c, id, updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Achievement updated"})
	})

	// Soft delete achievement (admin)
	a.DELETE(":id", middleware.RBAC("admin"), func(c *gin.Context) {
		id := c.Param("id")
		if err := achievementService.SoftDelete(c, id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Achievement deleted"})
	})

	// Add attachment to achievement
	a.POST(":id/attachments", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		id := c.Param("id")
		var att models.Attachment
		if err := c.ShouldBindJSON(&att); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if err := achievementService.AddAttachment(c, id, att); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Attachment added"})
	})

	// Remove attachment by file name
	a.DELETE(":id/attachments/:filename", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		id := c.Param("id")
		filename := c.Param("filename")
		if err := achievementService.RemoveAttachment(c, id, filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Attachment removed"})
	})
}
