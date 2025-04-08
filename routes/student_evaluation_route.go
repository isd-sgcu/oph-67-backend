package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/handler"
	"github.com/isd-sgcu/oph-67-backend/middleware"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

func RegisterStudentEvaluationRoutes(app *fiber.App, studentEvaluationUsecase *usecase.StudentEvaluationUsecase, userUsecases *usecase.UserUsecase) {
	studentEvaluationHandler := handler.NewStudentEvaluationHandler(studentEvaluationUsecase)

	api := app.Group("/api")

	// Authenticated user routes - Requires valid JWT
	authenticated := api.Group("/student-evaluation", middleware.AuthMiddleware(userUsecases))
	authenticated.Post("/", studentEvaluationHandler.CreateStudentEvaluation)           // Create a new student evaluation
	authenticated.Get("/:id", studentEvaluationHandler.GetStudentEvaluationByStudentId) // Get student evaluation by ID
	authenticated.Patch("/:id", studentEvaluationHandler.UpdateStudentEvaluation)       // Update student evaluation
	authenticated.Delete("/:id", studentEvaluationHandler.DeleteStudentEvaluation)      // Delete student evaluation

	// Staff/Admin routes - Requires Staff or Admin role
	staffAdmin := api.Group("/student-evaluation", middleware.RoleMiddleware(userUsecases, domain.Staff, domain.Admin))
	staffAdmin.Get("/", studentEvaluationHandler.GetAllStudentEvaluations) // List all student evaluations
}
