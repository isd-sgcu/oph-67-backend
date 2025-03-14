package handler

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/usecase"
	"github.com/isd-sgcu/oph-67-backend/utils"
	"github.com/lib/pq"
)

// UserHandler represents the handler for user-related endpoints
type UserHandler struct {
	Usecase *usecase.UserUsecase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: usecase}
}

// Register Staff godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Accept  multipart/form-data
// @Produce  json
// @Param id formData string true "ID"
// @Param name formData string true "Name"
// @Param phone formData string true "Phone"
// @Param email formData string true "Email"
// @Param faculty formData string true "Faculty"
// @Param Year	formData  int  true "true"
// @Param Nickname formData string true "Nickname"
// @Param StudentID formData string true "StudentID"
// @Param IsCentralStaff formData boolean true "IsCentralStaff"
// @Success 201 {object} domain.TokenResponse
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 500 {object} domain.ErrorResponse "Failed to create user"
// @Router /api/staff/register [post]
func (h *UserHandler) StaffRegister(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}

	// Helper function to get required form values
	getFormValue := func(key string) (string, bool) {
		if v, ok := form.Value[key]; ok && len(v) > 0 {
			return v[0], true
		}
		return "", false
	}

	// Helper function for optional values
	getOptionalValue := func(key string) *string {
		if v, ok := form.Value[key]; ok && len(v) > 0 {
			return &v[0]
		}
		return nil
	}

	// Initialize user data
	userData := make(map[string]string)

	// Validate phone number
	if phone, ok := getFormValue("phone"); ok {
		if !utils.IsValidPhone(phone) {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid phone number format"})
		}
		userData["phone"] = phone
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Phone is required"})
	}

	// Validate required fields
	requiredFields := []string{"id", "name", "phone", "email"}
	for _, field := range requiredFields {
		if value, ok := getFormValue(field); ok {
			userData[field] = value
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: field + " is required"})
		}
	}

	user := &domain.User{
		ID:        userData["id"],
		Name:      userData["name"],
		Nickname:  getOptionalValue("nickname"),
		StudentID: getOptionalValue("studentId"),
		Role:      domain.Staff,
		Email:     userData["email"],
		Phone:     userData["phone"],
		Faculty:   getOptionalValue("faculty"),
		Year: func() *int {
			if yearStr := getOptionalValue("year"); yearStr != nil {
				if year, err := strconv.Atoi(*yearStr); err == nil {
					return &year
				}
			}
			return nil
		}(),
		IsCentralStaff: func() *bool {
			if isCentralStaffStr := getOptionalValue("isCentralStaff"); isCentralStaffStr != nil {
				if isCentralStaff, err := strconv.ParseBool(*isCentralStaffStr); err == nil {
					return &isCentralStaff
				}
			}
			return nil
		}(),
	}

	tokenResponse, err := h.Usecase.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(tokenResponse)
}

// Register Student godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Accept  multipart/form-data
// @Produce  json
// @Param id formData string true "ID"
// @Param name formData string true "Name"
// @Param phone formData string true "Phone"
// @Param email formData string true "Email"
// @Param status formData string true "Status"
// @Param otherStatus formData string false "OtherStatus"
// @Param province formData string true "Province"
// @Param school formData string true "School"
// @Param selectedSources formData string true "SelectedSources"
// @Param birthDate formData string false "BirthDate"
// @Param otherSource formData string false "OtherSource"
// @Param firstInterest formData string true "FirstInterest"
// @Param secondInterest formData string true "SecondInterest"
// @Param thirdInterest formData string true "ThirdInterest"
// @Param objective formData string true "Objective"
// @Success 201 {object} domain.TokenResponse
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 500 {object} domain.ErrorResponse "Failed to create user"
// @Router /api/student/register [post]
func (h *UserHandler) StudentRegister(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}

	// Helper function to get required form values
	getFormValue := func(key string) (string, bool) {
		if v, ok := form.Value[key]; ok && len(v) > 0 {
			return v[0], true
		}
		return "", false
	}

	// Helper function for optional values
	getOptionalValue := func(key string) *string {
		if v, ok := form.Value[key]; ok && len(v) > 0 {
			return &v[0]
		}
		return nil
	}

	// Initialize user data
	userData := make(map[string]string)

	// Validate phone number
	if phone, ok := getFormValue("phone"); ok {
		if !utils.IsValidPhone(phone) {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid phone number format"})
		}
		userData["phone"] = phone
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Phone is required"})
	}

	// Validate required fields
	requiredFields := []string{"id", "name", "phone", "email"}
	for _, field := range requiredFields {
		if value, ok := getFormValue(field); ok {
			userData[field] = value
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: field + " is required"})
		}
	}

	selectedSources := getOptionalValue("selectedSources")
	var selectedSourcesArray *pq.StringArray

	if selectedSources != nil && *selectedSources != "" {
		arr := pq.StringArray(strings.Split(*selectedSources, ","))
		selectedSourcesArray = &arr
	} else {
		selectedSourcesArray = nil // or &pq.StringArray{} if you prefer an empty array
	}

	user := &domain.User{
		ID:    userData["id"],
		Name:  userData["name"],
		Role:  domain.Student,
		Email: userData["email"],
		Phone: userData["phone"],
		BirthDate: func() *time.Time {
			if birthDate := getOptionalValue("birthDate"); birthDate != nil {
				if parsedDate, err := time.Parse("2006-01-02", *birthDate); err == nil {
					return &parsedDate
				}
			}
			return nil
		}(),

		Status:          getOptionalValue("status"),
		OtherStatus:     getOptionalValue("otherStatus"),
		Province:        getOptionalValue("province"),
		School:          getOptionalValue("school"),
		SelectedSources: selectedSourcesArray,
		OtherSource:     getOptionalValue("otherSource"),
		FirstInterest:   getOptionalValue("firstInterest"),
		SecondInterest:  getOptionalValue("secondInterest"),
		ThirdInterest:   getOptionalValue("thirdInterest"),
		Objective:       getOptionalValue("objective"),
	}

	tokenResponse, err := h.Usecase.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(tokenResponse)
}

// GetAll godoc
// @Summary Get all student
// @Description Retrieve a list of all users with optional filtering
// @Produce  json
// @Param name query string false "Filter by name"
// @Param role query string false "Filter by role"
// @security BearerAuth
// @Success 200 {array} domain.User
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch users"
// @Router /api/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	// Get query parameters
	filter := c.Query("name")
	role := c.Query("role")

	users, err := h.Usecase.GetAll(filter, domain.Role(role))
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
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} domain.User
// @Failure 500 {object} domain.ErrorResponse "Failed to fetch User"
// @Failure 400 {object} domain.ErrorResponse "User has already entered"
// @Router /api/users/qr/{id} [post]
func (h *UserHandler) ScanQR(c *fiber.Ctx) error {
	// Extract student ID from URL params
	studentId := c.Params("id")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Missing token"})
	}

	// Extract staff ID from JWT token
	staffId, err := utils.DecodeToken(strings.TrimPrefix(authHeader, "Bearer "), utils.GetEnv("SECRET_JWT_KEY", ""))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Unauthorized"})
	}

	// Call use case with both student and staff IDs
	user, err := h.Usecase.ScanQR(studentId, staffId)
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
// @Param role body domain.RoleRequest true "User Role"
// @Success 204
// @Failure 400 {object} domain.ErrorResponse "Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to update user role"
// @Router /api/users/role/{id} [patch]
func (h *UserHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	role := new(domain.RoleRequest)
	if err := c.BodyParser(role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}
	if err := h.Usecase.UpdateRole(id, domain.Role(role.Role)); err != nil {
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

// RemoveStaff godoc
// @Summary RemoveStaff user by ID
// @Description RemoveStaff a user by its ID
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to delete user"
// @Router /api/users/{id} [delete]
func (h *UserHandler) RemoveStaff(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.Usecase.RemoveStaff(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed to delete user"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Delete user godoc
// @Summary RemoveStaff user by ID
// @Description RemoveStaff a user by its ID
// @Produce  json
// @security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Failure 404 {object} domain.ErrorResponse "User not found"
// @Failure 500 {object} domain.ErrorResponse "Failed to delete user"
// @Router /api/admin/users/{id} [delete]
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
	// recieve id from body in plain text
	id := new(domain.SignInRequest)
	if err := c.BodyParser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid input"})
	}

	tokenResponse, err := h.Usecase.SignIn(id.ID)
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
