package routes

import (
	"net/http"

	"backend/middleware"
	"backend/app/models"
	"backend/app/services"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(r *gin.Engine, studentService service.StudentService) {

	students := r.Group("/students")
	students.Use(middleware.JWTMiddleware()) // semua endpoint harus login

	students.GET("/", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {

		data, err := studentService.FindAll(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   data,
		})
	})

	students.GET("/:id", middleware.RBAC("admin", "dosen"), func(c *gin.Context) {
		id := c.Param("id")

		data, err := studentService.FindByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   data,
		})
	})

	students.POST("/", middleware.RBAC("admin"), func(c *gin.Context) {

		var req models.Student

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		err := studentService.Create(c, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Student created successfully",
		})
	})

	students.PUT("/:id", middleware.RBAC("admin"), func(c *gin.Context) {
		id := c.Param("id")

		var req models.Student

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		err := studentService.Update(c, id, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Student updated successfully",
		})
	})

	students.DELETE("/:id", middleware.RBAC("admin"), func(c *gin.Context) {
		id := c.Param("id")

		err := studentService.Delete(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Student deleted successfully",
		})
	})
}
