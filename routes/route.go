package routes

import (
	"os"

	"backend/app/models"
	"backend/app/service"
	"backend/middleware"

	_ "backend/docs"

	"github.com/gofiber/swagger"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService service.UserService, authService service.AuthService, achievementService service.AchievementService,
	referenceService service.AchievementReferenceService, studentLecturerService service.StudentLecturerService, reportService service.ReportService) {

	app.Get("/swagger/*", swagger.HandlerDefault)
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/login", processLogin(authService))
	auth.Post("/refresh", processRefreshToken(authService))
	auth.Post("/logout", processLogout(authService))
	auth.Get("/profile", middleware.JWTMiddleware(), processGetProfile(userService))

	users := api.Group("/users")
	users.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_SECRET")),
		},
		Claims: &models.JWTClaims{},
	}))

	users.Use(middleware.OnlyAdmin())
	users.Get("/", processGetAllUsers(userService))
	users.Post("/", processCreateUser(userService))
	users.Get("/:id", processGetUserByID(userService))
	users.Put("/:id", processUpdateUser(userService))
	users.Delete("/:id", processDeleteUser(userService))
	users.Put("/:id/role", processUpdateUserRole(userService))

	achievements := api.Group("/achievements")
	achievements.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_SECRET")),
		},
		Claims: &models.JWTClaims{},
	}))

	achievements.Use(middleware.OnlyMahasiswa())
	achievements.Get("/", listAchievements(referenceService))
	achievements.Get("/:id", getAchievementDetail(achievementService))
	achievements.Post("/", createAchievement(achievementService, referenceService))
	achievements.Put("/:id", updateAchievement(achievementService))
	achievements.Delete("/:id", deleteAchievement(referenceService))
	achievements.Post("/:id/submit", submitAchievement(referenceService))
	achievements.Post("/:id/attachments", addAttachment(achievementService))

	achievements.Use(middleware.OnlyDosenWali())
	achievements.Post("/:id/verify", verifyAchievement(referenceService))
	achievements.Post("/:id/reject", rejectAchievement(referenceService))
	achievements.Get("/:id/history", achievementHistory(referenceService))

	students := api.Group("/students")
	students.Get("/", StudentList(studentLecturerService))
	students.Get("/:id", StudentGetByID(studentLecturerService))
	students.Put("/:id/advisor", StudentUpdateAdvisor(studentLecturerService))

	lecturers := api.Group("/lecturers")
	lecturers.Get("/", LecturerList(studentLecturerService))
	lecturers.Get("/:id/advisees", LecturerAdvisees(studentLecturerService))

	reports := api.Group("/reports", middleware.JWTMiddleware())
	reports.Get("/statistics", handleGetStatistics(reportService))
	reports.Get("/student/:id", handleGetStudentReport(reportService))

}
