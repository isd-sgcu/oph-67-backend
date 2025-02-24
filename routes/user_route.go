package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/handler"
	"github.com/isd-sgcu/oph-67-backend/middleware"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

// RegisterUserRoutes sets up all user-related endpoints with appropriate middleware and grouping
func RegisterUserRoutes(app *fiber.App, userUsecase *usecase.UserUsecase) {
	userHandler := handler.NewUserHandler(userUsecase)

	api := app.Group("/api/users")

	// Public routes - No authentication required
	{
		api.Post("/signin", userHandler.SignIn)     // User authentication
		api.Post("/student/register", userHandler.StudentRegister) // New student registration
	}

	// Authenticated user routes - Requires valid JWT
	authenticated := api.Group("", middleware.AuthMiddleware(userUsecase))
	{
		authenticated.Get("/:id", userHandler.GetById)            // Get user by ID (self)
		authenticated.Patch("/", userHandler.UpdateMyAccountInfo) // Update own account info
		authenticated.Get("/qr/:id", userHandler.GetQRURL)        // Get user's QR code URL
	}

	// Staff/Admin routes - Requires Staff or Admin role
	staffAdmin := api.Group("", middleware.RoleMiddleware(userUsecase, domain.Staff, domain.Admin))
	{
		staffAdmin.Get("/", userHandler.GetAll)        // List all users
		staffAdmin.Post("/qr/:id", userHandler.ScanQR) // Scan user QR code
	}

	// Admin-only routes - Requires Admin role
	admin := api.Group("", middleware.RoleMiddleware(userUsecase, domain.Admin))
	{
		admin.Patch("/:id", userHandler.Update)               // Update any user account
		admin.Delete("/:id", userHandler.Delete)              // Delete user
		admin.Patch("/role/:id", userHandler.UpdateRole)      // Update user role
		admin.Patch("/addstaff/:phone", userHandler.AddStaff) // Promote user to Staff by phone
	}
}
