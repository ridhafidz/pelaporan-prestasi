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

	achievements.Get("/", middleware.OnlyMahasiswa(), listAchievements(referenceService))
	achievements.Get("/:id", middleware.OnlyMahasiswa(), getAchievementDetail(achievementService))
	achievements.Post("/", middleware.OnlyMahasiswa(), createAchievement(achievementService, referenceService))
	achievements.Put("/:id", middleware.OnlyMahasiswa(), updateAchievement(achievementService))
	achievements.Delete("/:id", middleware.OnlyMahasiswa(), deleteAchievement(referenceService))
	achievements.Post("/:id/submit", middleware.OnlyMahasiswa(), submitAchievement(referenceService))
	achievements.Post("/:id/attachments", middleware.OnlyMahasiswa(), addAttachment(achievementService))

	achievements.Post("/:id/verify", middleware.OnlyDosenWali(), verifyAchievement(referenceService))
	achievements.Post("/:id/reject", middleware.OnlyDosenWali(), rejectAchievement(referenceService))
	achievements.Get("/:id/history", middleware.OnlyDosenWali(), achievementHistory(referenceService))

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
