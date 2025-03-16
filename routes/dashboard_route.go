package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/handler"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

func RegisterDashboardRoutes(app *fiber.App, dashboardUsecase *usecase.DashboardUseCase) {
	dashboardHandler := handler.NewDashBoardUseCase(dashboardUsecase)

	api := app.Group("/api")

	authenticated := api.Group("/dashboard")
	authenticated.Get("/faculties", dashboardHandler.GetFacultyCount)
	authenticated.Get("/sources", dashboardHandler.GetSourceCount)
}
