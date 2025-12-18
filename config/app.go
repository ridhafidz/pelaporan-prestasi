package config

import (
	"fmt"
	"log"
	"os"

	"backend/app/repository"
	"backend/app/service"
	"backend/database"
	"backend/routes"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	App *fiber.App
}

func InitApp() *Application {
	LoadENV()
	fmt.Println("SECRET KEY:", os.Getenv("JWT_SECRET"))

	database.ConnectPostgres()
	database.ConnectMongo()

	postgresDB := database.PostgresDB
	mongoDB := database.MongoDB

	userRepo := repository.NewUserRepository(postgresDB)
	authRepo := repository.NewAuthRepository(postgresDB)

	achievementRepo := repository.NewAchievementRepository(mongoDB)
	achievementRefRepo := repository.NewAchievementReferenceRepository(postgresDB)

	studentLecturerRepo := repository.NewStudentLecturerRepository(postgresDB)

	reportRepo := repository.NewReportRepository(postgresDB, mongoDB)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set")
	}

	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)

	achievementService := service.NewAchievementService(achievementRepo)
	achievementReferenceService :=
		service.NewAchievementReferenceService(
			achievementRefRepo,
			achievementRepo,
		)
	studentLecturerService := service.NewStudentLecturerService(studentLecturerRepo)
	reportService := service.NewReportService(reportRepo, studentLecturerRepo)

	app := fiber.New(fiber.Config{
		AppName: "Pelaporan-Prestasi",
	})

	routes.SetupRoutes(
		app,
		userService,
		authService,
		achievementService,
		achievementReferenceService,
		studentLecturerService,
		reportService,
	)

	log.Println("🚀 Application running on port:", os.Getenv("APP_PORT"))

	return &Application{App: app}
}
