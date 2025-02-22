package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/handler"
	"github.com/isd-sgcu/oph-67-backend/middleware"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

func RegisterUserRoutes(app *fiber.App, userUsecase *usecase.UserUsecase) {
	userHandler := handler.NewUserHandler(userUsecase)

	api := app.Group("/api/users")

	api.Post("/signin", userHandler.SignIn)

	api.Get("/",
		middleware.RoleMiddleware(
			userUsecase,
			domain.Staff,
			domain.Admin,
		),
		userHandler.GetAll)

	api.Get("/:id", middleware.AuthMiddleware(userUsecase), userHandler.GetById)

	api.Post("/qr/:id", middleware.RoleMiddleware(
		userUsecase,
		domain.Staff,
		domain.Admin,
	),
		userHandler.ScanQR)

	api.Get("/qr/:id", middleware.AuthMiddleware(userUsecase), userHandler.GetQRURL)

	api.Post("/register", userHandler.Register)

	api.Patch("/:id", middleware.RoleMiddleware(
		userUsecase,
		domain.Admin,
	),
		userHandler.Update)

	api.Patch("/", middleware.AuthMiddleware(userUsecase), userHandler.UpdateMyAccountInfo)
	api.Patch("/addstaff/:phone", middleware.RoleMiddleware(userUsecase, domain.Admin), userHandler.AddStaff)
	api.Delete("/:id", middleware.RoleMiddleware(userUsecase, domain.Admin), userHandler.Delete)
	api.Patch("/role/:id", middleware.RoleMiddleware(userUsecase, domain.Admin), userHandler.UpdateRole)
}
