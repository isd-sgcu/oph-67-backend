package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/usecase"
	"github.com/isd-sgcu/oph-67-backend/utils"
)

type StudentEvaluationHandler struct {
	Usecase *usecase.StudentEvaluationUsecase
}

func NewStudentEvaluationHandler(usecase *usecase.StudentEvaluationUsecase) *StudentEvaluationHandler {
	return &StudentEvaluationHandler{Usecase: usecase}
}

// CreateStudentEvaluation creates a new student evaluation.
func (h *StudentEvaluationHandler) CreateStudentEvaluation(c *fiber.Ctx) error {
	var evaluation domain.StudentEvaluation
	// Get student ID from authenticated user
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Missing token"})
	}
	userId, err := utils.DecodeToken(strings.TrimPrefix(authHeader, "Bearer "), utils.GetEnv("SECRET_JWT_KEY", ""))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid token"})
	}
	if err := c.BodyParser(&evaluation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err = h.Usecase.CreateStudentEvaluation(&domain.StudentEvaluation{
		StudentId:                         userId,
		NewSources:                        evaluation.NewSources,
		OverallActivity:                   evaluation.OverallActivity,
		InterestActivity:                  evaluation.InterestActivity,
		ReceivedFacultyInfoClearly:        evaluation.ReceivedFacultyInfoClearly,
		WouldRecommendCUOpenHouseNextTime: evaluation.WouldRecommendCUOpenHouseNextTime,
		FavoriteBooth:                     evaluation.FavoriteBooth,
		ActivityDiversity:                 evaluation.ActivityDiversity,
		PerceivedCrowdDensity:             evaluation.PerceivedCrowdDensity,
		HasFullBoothAccess:                evaluation.HasFullBoothAccess,
		FacilityConvenienceRating:         evaluation.FacilityConvenienceRating,
		CampusNavigationRating:            evaluation.CampusNavigationRating,
		HesitationLevelAfterDisaster:      evaluation.HesitationLevelAfterDisaster,
		LineOASignupRating:                evaluation.LineOASignupRating,
		DesignBeautyRating:                evaluation.DesignBeautyRating,
		WebsiteImprovementSuggestions:     evaluation.WebsiteImprovementSuggestions,
	})
	if err != nil {
		if err == domain.ErrStudentEvaluationAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Student evaluation already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(evaluation)
}

// GetStudentEvaluationByStudentId retrieves a student evaluation by student ID.
func (h *StudentEvaluationHandler) GetStudentEvaluationByStudentId(c *fiber.Ctx) error {
	studentId := c.Params("student_id")
	if studentId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Student ID is required"})
	}

	evaluation, err := h.Usecase.GetStudentEvaluationByStudentId(studentId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student evaluation not found"})
	}

	return c.JSON(evaluation)
}

// GetAllStudentEvaluations retrieves all student evaluations.
func (h *StudentEvaluationHandler) GetAllStudentEvaluations(c *fiber.Ctx) error {
	evaluations, err := h.Usecase.GetAllStudentEvaluations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(evaluations)
}

// UpdateStudentEvaluation updates an existing student evaluation.
func (h *StudentEvaluationHandler) UpdateStudentEvaluation(c *fiber.Ctx) error {
	studentId := c.Params("student_id")
	if studentId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Student ID is required"})
	}

	var evaluation domain.StudentEvaluation
	if err := c.BodyParser(&evaluation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	evaluation.StudentId = studentId

	err := h.Usecase.UpdateStudentEvaluation(&evaluation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(evaluation)
}

// DeleteStudentEvaluation deletes a student evaluation by student ID.
func (h *StudentEvaluationHandler) DeleteStudentEvaluation(c *fiber.Ctx) error {
	studentId := c.Params("student_id")
	if studentId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Student ID is required"})
	}

	err := h.Usecase.DeleteStudentEvaluation(studentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
