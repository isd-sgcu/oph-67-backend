package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/handler"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

func RegisterDashboardRoutes(app *fiber.App, dashboardUsecase *usecase.DashboardUseCase) {
	dashboardHandler := handler.NewDashBoardUseCase(dashboardUsecase)

	api := app.Group("/api")

	dashboard := api.Group("/dashboard")
	dashboard.Get("/faculties", dashboardHandler.GetFacultyCount)
	dashboard.Get("/sources", dashboardHandler.GetSourceCount)
	dashboard.Get("/ages", dashboardHandler.GetAgeGroupCount)
	dashboard.Get("/faculties/today", dashboardHandler.GetFacultyTodayCount)
	dashboard.Get("/status", dashboardHandler.GetStatusStudent)
	dashboard.Get("/download", dashboardHandler.ExportAllStudents)
}
