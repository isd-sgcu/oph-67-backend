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
	Repo UserRepositoryInterface
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

// NewUserUsecase initializes a new UserUsecase instance with the provided repository.
func NewUserUsecase(repo UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

// assignRole determines and assigns a user's role based on their phone number.
// This is a temporary implementation using mock phone number lists for demonstration.
// TODO: Replace mock data with a proper role assignment mechanism (e.g., database lookup).
func (u *UserUsecase) assignRole(user *domain.User) {
	staffPhones := []string{"06", "05", "04", "03", "02"} // Mock staff phone prefixes
	adminPhones := []string{"01", "00", "99", "98", "97"} // Mock admin phone prefixes
	user.Role = domain.Member

	if user.Phone != "" {
		for _, phone := range staffPhones {
			if user.Phone == phone {
				user.Role = domain.Staff
				break
			}
		}

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

// Register handles user registration by generating a unique UID, assigning a role,
// persisting the user, and returning authentication tokens.
// Returns TokenResponse with access token or error if any step fails.
func (u *UserUsecase) Register(user *domain.User) (domain.TokenResponse, error) {
	u.assignRole(user)

	// Validate phone number
	if !utils.IsValidPhone(user.Phone) {
		return domain.TokenResponse{}, domain.ErrInvalidPhone
	}

	// Generate unique UID ensuring no collisions
	for {
		user.UID = utils.GenerateUID()
		uidExists, err := u.Repo.IsUIDExists(user.UID)
		if err != nil {
			return domain.TokenResponse{}, fmt.Errorf("error checking UID uniqueness: %w", err)
		}
		if !uidExists {
			break
		}
	}

	now := time.Now()
	user.RegisteredAt = &now

	// Persist user to database
	if err := u.Repo.Create(user); err != nil {
		return domain.TokenResponse{}, fmt.Errorf("error saving user: %w", err)
	}

	// Generate JWT tokens for authentication
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
func (u *UserUsecase) GetAll(filter string) ([]domain.User, error) {
	if filter != "" {
		users, err := u.Repo.GetByName(filter)
		if err != nil {
			return nil, err
		}
		return users, nil
	}

	return u.Repo.GetAll()
}

// GetById fetches a single user by their unique ID.
// Returns the user or error if not found or repository operation fails.
func (u *UserUsecase) GetById(id string) (domain.User, error) {
	return u.Repo.GetById(id)
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

	return u.Repo.Update(id, updatedUser)
}

// ScanQR records a user's entry by updating their LastEntered timestamp.
// Returns error if user has already entered today or repository operation fails.
func (u *UserUsecase) ScanQR(id string) (domain.User, error) {
	user, err := u.GetById(id)
	if err != nil {
		return domain.User{}, err
	}

	now := time.Now()
	if user.LastEntered != nil && isSameDay(*user.LastEntered, now) {
		return user, domain.ErrUserAlreadyEntered
	}

	user.LastEntered = &now
	err = u.Update(id, &user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
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

// Delete removes a user from the system by their ID.
// Returns error if repository operation fails.
func (u *UserUsecase) Delete(id string) error {
	return u.Repo.Delete(id)
}

// AddStaff promotes a user to staff role by looking up their phone number.
// Returns error if user not found, already staff, or update fails.
func (u *UserUsecase) AddStaff(phone string) error {
	user, err := u.Repo.GetByPhone(phone)
	if err != nil {
		return err
	}

	if user.Role == domain.Staff {
		return domain.ErrUserAlreadyStaff
	}

	user.Role = domain.Staff
	return u.Update(user.ID, &user)
}
