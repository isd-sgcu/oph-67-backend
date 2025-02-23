package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/usecase"
	"github.com/isd-sgcu/oph-67-backend/utils"
)

// UserHandler represents the handler for user-related endpoints
type UserHandler struct {
	Usecase *usecase.UserUsecase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: usecase}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Accept  multipart/form-data
// @Produce  json
// @Param user body domain.User true "User data"
// @Success 201 {object} domain.TokenResponse
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 500 {object} domain.ErrorResponse "Failed to create user"
// @Router /api/users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}

	tokenResponse, err := h.Usecase.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(tokenResponse)
}

// GetAll godoc
// @Summary Get all users
// @Description Retrieve a list of all users with optional filtering
// @Produce  json
// @Param name query string false "Filter by name"
// @security BearerAuth
// @Success 200 {array} domain.User
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch users"
// @Router /api/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	// Get query parameters
	filter := c.Query("name")

	users, err := h.Usecase.GetAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to fetch users"})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// GetById godoc
// @Summary Get user by ID
// @Description Retrieve a user by its ID
// @Accept  json
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch user"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.Usecase.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// Update godoc
// @Summary Update user by ID
// @Description Update a user by its ID
// @Accept  json
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Param user body domain.User true "User data"
// @Success 204
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to update user"
// @Router /api/users/{id} [patch]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}
	if err := h.Usecase.Update(id, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to update user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Scan QR godoc
// @Summary Scan QR code
// @Description Retrieve a user by its ID
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch User"
// @Failure 400 {object} domain.ErrorResponse "User has already entered"
// @Router /api/users/qr/{id} [post]
func (h *UserHandler) ScanQR(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.Usecase.ScanQR(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyEntered) {
			t := user.LastEntered.String()
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "User has already entered", Message: &t})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to scan QR"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// Change Role godoc
// @Summary Update user role by ID
// @Description Update a user by its ID
// @Accept  json
// @security BearerAuth
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param role body domain.Role true "User Role"
// @Success 204
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to update user role"
// @Router /api/users/role/{id} [patch]
func (h *UserHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	role := new(domain.Role)
	if err := c.BodyParser(role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}
	if err := h.Usecase.UpdateRole(id, *role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to update this role user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Update account info godoc
// @Summary Update Account Info
// @Description Update a user by its ID
// @Accept  json
// @Produce  json
// @security BearerAuth
// @Param user body domain.User true "User data"
// @Success 204
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to update user role"
// @Router /api/users [patch]
func (h *UserHandler) UpdateMyAccountInfo(c *fiber.Ctx) error {
	user := new(domain.User)
	tokenHeader := c.Get("Authorization")
	if !strings.HasPrefix(tokenHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Unauthorized"})
	}
	token := strings.TrimPrefix(tokenHeader, "Bearer ")

	id, err := utils.DecodeToken(token, utils.GetEnv("SECRET_JWT_KEY", ""))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Unauthorized " + err.Error()})
	}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}
	if err := h.Usecase.Update(id, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to update this role user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetQRURL godoc
// @Summary Get QR code URL
// @Description Retrieve a QR code URL for a user
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} domain.QrResponse
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch user"
// @Router /api/users/qr/{id} [get]
func (h *UserHandler) GetQRURL(c *fiber.Ctx) error {
	id := c.Params("id")
	qrURL, err := h.Usecase.GetQRURL(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(domain.QrResponse{QrURL: qrURL})
}

// Delete godoc
// @Summary Delete user by ID
// @Description Delete a user by its ID
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to delete user"
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.Usecase.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to delete user"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Produce  json
// @Param id body string true "User ID"
// @Success 200 {object} domain.TokenResponse
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 500 {object} domain.ErrorResponse "Failed to signin"
// @Router /api/users/signin [post]
func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	id := new(string)
	if err := c.BodyParser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}

	tokenResponse, err := h.Usecase.SignIn(*id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to signin"})
	}

	return c.Status(fiber.StatusOK).JSON(tokenResponse)
}

// Add Staff godoc
// @Summary Add Staff
// @security BearerAuth
// @Description Add Staff By phone number
// @Produce  json
// @Param phone path string true "User Phone"
// @Success 204
// @Failure 400 {object} domain.ErrorResponse "User is already a staff"
// @Failure 500 {object} domain.ErrorResponse "Failed to add staff"
// @Router /api/users/addstaff/{phone} [patch]
func (h *UserHandler) AddStaff(c *fiber.Ctx) error {
	phone := c.Params("phone")
	if err := h.Usecase.AddStaff(phone); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyStaff) {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "User is already a staff"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to add staff"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
