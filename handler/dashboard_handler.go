package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/usecase"
)

type DashBoardHandler struct {
	Usecase *usecase.DashboardUseCase
}

func NewDashBoardUseCase(usecase *usecase.DashboardUseCase) *DashBoardHandler {
	return &DashBoardHandler{Usecase: usecase}
}

// GetFacultyCount returns the number of students interested in each faculty.
func (h *DashBoardHandler) GetFacultyCount(c *fiber.Ctx) error {
	results, err := h.Usecase.GetFacultyCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

// GetSourceCount returns the number of students who selected each source.
func (h *DashBoardHandler) GetSourceCount(c *fiber.Ctx) error {
	results, err := h.Usecase.GetSourceCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}
