package config

import (
	"context"
	"log"
	"os"

	"backend/app/repository/postgree"
	"backend/app/repository/mongo"
	"backend/app/services/mongo"
	"backend/app/services/postgree"
	"backend/database"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Router *gin.Engine
}

func InitApp() *Application {
	LoadENV()

	dbx, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// sqlx.DB contains underlying *sql.DB as field DB
	stdDB := dbx.DB

	userRepo := repository.NewUserRepository(stdDB)
	authRepo := repository.NewAuthRepository(stdDB)
	studentRepo := repository.NewStudentRepository(dbx)

	// create mongo client for achievements repository
	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "prestasi"
	}
	achievementRepo := repository.NewAchievementRepository(mongoClient, dbName)

	authService := service.NewAuthService(userRepo, authRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)
	achievementService := service.NewAchievementService(achievementRepo)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.RegisterAuthRoutes(r, authService)
	routes.RegisterUserRoutes(r, userService)
	routes.RegisterStudentRoutes(r, studentService)
	routes.RegisterAchievementRoutes(r, achievementService)

	log.Println("Application started on port:", ENV.AppPort)

	return &Application{
		Router: r,
	}
}
