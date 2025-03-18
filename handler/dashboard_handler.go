package handler

import (
	"compress/gzip"
	"encoding/csv"

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

// GetAgeGroupCount returns the number of students in each age group.
func (h *DashBoardHandler) GetAgeGroupCount(c *fiber.Ctx) error {
	results, err := h.Usecase.GetAgeGroupCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

// GetFacultyTodayCount returns the number of students interested in each faculty today.
func (h *DashBoardHandler) GetFacultyTodayCount(c *fiber.Ctx) error {
	results, err := h.Usecase.GetFacultyTodayCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

// GetStatusStudent returns the number of students in each status.
func (h *DashBoardHandler) GetStatusStudent(c *fiber.Ctx) error {
	results, err := h.Usecase.GetStatusStudent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func (h *DashBoardHandler) ExportAllStudents(c *fiber.Ctx) error {
	students, err := h.Usecase.GetAllStudent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch student data",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=students_export.csv")
	c.Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(c.Response().BodyWriter())
	defer gz.Close()
	writer := csv.NewWriter(gz)

	defer writer.Flush()

	header := []string{
		"ID",
		"Name",
		"Email",
		"Phone",
		"First Interest",
		"Second Interest",
		"Third Interest",
	}
	if err := writer.Write(header); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate CSV",
		})
	}

	// เขียนข้อมูลทีละแถว
	for _, student := range students {
		row := []string{
			student.ID,
			student.Name,
			student.Email,
			student.Phone,
			getSafeString(&student.FirstInterest),
			getSafeString(&student.SecondInterest),
			getSafeString(&student.ThirdInterest),
		}
		if err := writer.Write(row); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to write CSV row",
			})
		}
	}

	return nil
}

// Helper function สำหรับจัดการค่า Null
func getSafeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
