package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/isd-sgcu/oph-67-backend/config"
	_ "github.com/isd-sgcu/oph-67-backend/docs"
	"github.com/isd-sgcu/oph-67-backend/infrastructure"
	"github.com/isd-sgcu/oph-67-backend/middleware"
	"github.com/isd-sgcu/oph-67-backend/repository"
	"github.com/isd-sgcu/oph-67-backend/routes"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Fiber app
	app := fiber.New()

	// Add middleware
	app.Use(middleware.RequestLoggerMiddleware())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           // Allowed origin
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",           // Allow all necessary HTTP methods
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Include Authorization and other headers
	}))

	// Connect to the database
	db := infrastructure.ConnectDatabase(cfg)

	// Connect to Cache

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	dashBoardRepo := repository.NewDashBoardRepository(db)
	transactionRepo := repository.NewStudentTransactionRepository(db)
	studentEvaluationRepo := repository.NewStudentEvaluationRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, transactionRepo)
	dashBoardUssecase := usecase.NewDashBoardUseCase(dashBoardRepo)
	studentEvaluationUsecase := usecase.NewStudentEvaluationUsecase(studentEvaluationRepo)

	// Register routes
	routes.RegisterUserRoutes(app, userUsecase) // Register the user routes
	routes.RegisterDashboardRoutes(app, dashBoardUssecase)
	routes.RegisterStudentEvaluationRoutes(app, studentEvaluationUsecase, userUsecase)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/doc.json", // URL to access the Swagger docs
	}))

	// Start the server
	if err := app.Listen(":4000"); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
