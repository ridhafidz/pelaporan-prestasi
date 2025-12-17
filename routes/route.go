package routes

import (
	"os"

	"backend/app/services"
	"backend/middleware"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService service.UserService, authService service.AuthService) {

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

}
