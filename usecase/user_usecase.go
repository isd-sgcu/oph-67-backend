package usecase

import (
	"fmt"
	"time"

	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/utils"
)

// UserUsecase provides business logic operations for user management.
// It depends on a UserRepositoryInterface to interact with the data layer.
type UserUsecase struct {
	UserRepo               UserRepositoryInterface
	StudentTransactionRepo StudentTransactionRepositoryInterface
}

// UserRepositoryInterface defines the repository methods required by UserUsecase.
// Implementations of this interface handle data storage and retrieval operations.
type UserRepositoryInterface interface {
	Create(user *domain.User) error
	GetAll() ([]domain.User, error)
	GetById(id string) (domain.User, error)
	GetByPhone(phone string) (domain.User, error)
	GetByName(name string) ([]domain.User, error)
	IsUIDExists(uid string) (bool, error)
	Update(id string, user *domain.User) error
	Delete(id string) error
}

type StudentTransactionRepositoryInterface interface {
	Create(transaction *domain.StudentTransaction) error
	GetAll() ([]domain.StudentTransaction, error)
	GetById(id string) (domain.StudentTransaction, error)
	GetByStudentId(studentId string) ([]domain.StudentTransaction, error)
	GetByStudentIdAndFaculty(studentId string, faculty string) ([]domain.StudentTransaction, error)
	Update(id string, transaction *domain.StudentTransaction) error
	Delete(id string) error
}

// NewUserUsecase initializes a new UserUsecase instance with the provided repository.
func NewUserUsecase(userRepo UserRepositoryInterface, studentTransactionRepo StudentTransactionRepositoryInterface) *UserUsecase {
	return &UserUsecase{UserRepo: userRepo, StudentTransactionRepo: studentTransactionRepo}
}

// assignRole determines and assigns a user's role based on their phone number.
// This is a temporary implementation using mock phone number lists for demonstration.
// TODO: Replace mock data with a proper role assignment mechanism (e.g., database lookup).
func (u *UserUsecase) assignRole(user *domain.User) {
	adminPhones := []string{"0949823195"} // Mock admin phone prefixes

	if user.Phone != "" {
		for _, phone := range adminPhones {
			if user.Phone == phone {
				user.Role = domain.Admin
				break
			}
		}
	}
}

// isSameDay checks if two time values occur on the same calendar day.
// It compares year, month, and day components while ignoring time.
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (u *UserUsecase) Register(user *domain.User) (domain.TokenResponse, error) {
	u.assignRole(user)

	// Ensure UID is unique with a loop limit
	maxAttempts := 10
	for attempts := 0; attempts < maxAttempts; attempts++ {
		user.UID = utils.GenerateUID()
		uidExists, err := u.UserRepo.IsUIDExists(user.UID)
		if err != nil {
			return domain.TokenResponse{}, fmt.Errorf("error checking UID uniqueness: %w", err)
		}
		if !uidExists {
			break
		}
	}

	now := time.Now()
	user.RegisteredAt = &now

	// Check if user already exists
	existingUser, err := u.UserRepo.GetById(user.ID)
	if err != nil {
		fmt.Println("Error fetching user", err)
		fmt.Println("Trying to create user")
		// User not found, create a new one
		if err := u.UserRepo.Create(user); err != nil {
			return domain.TokenResponse{}, fmt.Errorf("error saving user: %w", err)
		}
		return u.generateTokenResponse(user)
	}

	// Promote to staff if already exists and is a member
	if existingUser.Role == domain.Member {
		existingUser.Role = domain.Staff
		if err := u.UserRepo.Update(existingUser.ID, &existingUser); err != nil {
			return domain.TokenResponse{}, fmt.Errorf("error updating user role: %w", err)
		}
	}

	return u.generateTokenResponse(&existingUser)
}

func (u *UserUsecase) generateTokenResponse(user *domain.User) (domain.TokenResponse, error) {
	jwtSecret := utils.GetEnv("SECRET_JWT_KEY", "")
	accessToken, err := utils.GenerateTokens(user.ID, jwtSecret)
	if err != nil {
		return domain.TokenResponse{}, fmt.Errorf("error generating tokens: %w", err)
	}

	return domain.TokenResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	}, nil
}

// GetAll retrieves users, optionally filtered by name if the filter parameter is provided.
// Returns a list of users or error if repository operation fails.
func (u *UserUsecase) GetAll(filter string, role domain.Role) ([]domain.User, error) {
	if filter != "" || role != "" {
		users, err := u.UserRepo.GetByName(filter)
		if role != "" {
			var filteredUsers []domain.User
			for _, user := range users {
				if user.Role == role {
					filteredUsers = append(filteredUsers, user)
				}
			}
			return filteredUsers, nil
		}
		if err != nil {
			return nil, err
		}
		return users, nil
	}

	return u.UserRepo.GetAll()
}

// GetById fetches a single user by their unique ID.
// Returns the user or error if not found or repository operation fails.
func (u *UserUsecase) GetById(id string) (domain.User, error) {
	return u.UserRepo.GetById(id)
}

// SignIn generates new authentication tokens for an existing user.
// Returns TokenResponse with access token or error if user lookup fails.
func (u *UserUsecase) SignIn(id string) (domain.TokenResponse, error) {
	user, err := u.GetById(id)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	jwtSecret := utils.GetEnv("SECRET_JWT_KEY", "")
	accessToken, err := utils.GenerateTokens(user.ID, jwtSecret)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	}, nil
}

// Update modifies an existing user's information.
// Returns error if user doesn't exist or repository operation fails.
func (u *UserUsecase) Update(id string, updatedUser *domain.User) error {
	_, err := u.GetById(id)
	if err != nil {
		return err
	}

	return u.UserRepo.Update(id, updatedUser)
}

// ScanQR records a user's entry by updating their LastEntered timestamp.
// Returns error if user has already entered today or repository operation fails.
func (u *UserUsecase) ScanQR(studentId string, staffId string) (domain.User, error) {
	student, err := u.GetById(studentId)
	if err != nil {
		return domain.User{}, err
	}

	staff, err := u.GetById(staffId)
	if err != nil {
		return domain.User{}, err
	}

	now := time.Now()

	if staff.Faculty != nil {
		return u.processFacultyStaffEntry(studentId, *staff.Faculty, now, student)
	}

	return u.processCentralStaffEntry(studentId, &student, now)
}

// UpdateRole changes a user's role to the specified value.
// Typically used by administrators for role management.
// Returns error if user doesn't exist or update fails.
func (u *UserUsecase) UpdateRole(id string, role domain.Role) error {
	user, err := u.GetById(id)
	if err != nil {
		return err
	}
	user.Role = role
	return u.Update(id, &user)
}

// GetQRURL generates the full URL for a user's QR code based on their ID.
// Uses the PRODUCTION_BASE_URL environment variable to construct the URL.
func (u *UserUsecase) GetQRURL(id string) (string, error) {
	user, err := u.GetById(id)
	if err != nil {
		return "", err
	}

	baseURL := utils.GetEnv("PRODUCTION_BASE_URL", "http://localhost:4000")

	return fmt.Sprintf("%s/api/users/qr/%s", baseURL, user.ID), nil
}

// RemoveStaff removes a user from the system by their ID.
// Returns error if repository operation fails.
func (u *UserUsecase) RemoveStaff(id string) error {
	return u.UserRepo.Update(id, &domain.User{Role: domain.Member})
}

// AddStaff promotes a user to staff role by looking up their phone number.
// Returns error if user not found, already staff, or update fails.
func (u *UserUsecase) AddStaff(phone string) error {
	user, err := u.UserRepo.GetByPhone(phone)
	if err != nil {
		return err
	}

	if user.Role == domain.Staff {
		return domain.ErrUserAlreadyStaff
	}

	user.Role = domain.Staff
	return u.Update(user.ID, &user)
}

// Delete All user
func (u *UserUsecase) Delete(id string) error {
	return u.UserRepo.Delete(id)
}

func (u *UserUsecase) isNull(staff domain.User) bool {
	return staff.IsCentralStaff == nil
}

// func (u *UserUsecase) hasEnteredToday(lastEntered *time.Time, now time.Time) bool {
// 	return lastEntered != nil && isSameDay(*lastEntered, now)
// }

func (u *UserUsecase) processCentralStaffEntry(studentId string, student *domain.User, now time.Time) (domain.User, error) {
	student.LastEntered = &now
	student.Faculty = nil
	err := u.Update(studentId, student)
	if err != nil {
		return domain.User{}, err
	}
	return *student, nil
}

func (u *UserUsecase) processFacultyStaffEntry(studentId, faculty string, now time.Time, student domain.User) (domain.User, error) {
	existingTransactions, err := u.StudentTransactionRepo.GetByStudentIdAndFaculty(studentId, faculty)
	if err != nil {
		return domain.User{}, err
	}

	for _, transaction := range existingTransactions {
		if isSameDay(transaction.RegisteredAt, now) {
			return student, domain.ErrUserAlreadyEntered
		}
	}

	err = u.StudentTransactionRepo.Create(&domain.StudentTransaction{
		ID:                    utils.GenerateUID(),
		StudentRegistrationID: studentId,
		Faculty:               faculty,
		RegisteredAt:          now,
	})
	if err != nil {
		return domain.User{}, err
	}

	return student, nil
}
