package routes

import (
	"os"

	"backend/app/service"
	"backend/middleware"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService service.UserService, authService service.AuthService, achievementService service.AchievementService,
	referenceService service.AchievementReferenceService) {

	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/login", processLogin(authService))
	auth.Post("/refresh", processRefreshToken(authService))
	auth.Post("/logout", processLogout(authService))
	auth.Get("/profile", processGetProfile(userService))

	users := api.Group("/users")
	users.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
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
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
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

}
