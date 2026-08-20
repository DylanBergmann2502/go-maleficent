// internal/app/users/services/user_service.go
package services

import (
	stderrors "errors"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/utils"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = apperrors.NewApplicationError(
		apperrors.ErrConflictCategory,
		"user_already_exists",
		"user with this email already exists",
		nil,
	)
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser handles the business logic of registering a new user
func (s *UserService) CreateUser(email, password string) (*models.User, error) {
	// 1. Check if user exists
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, apperrors.WrapApplicationError(
			err,
			apperrors.ErrInternalCategory,
			"database_error",
			"failed to check whether user exists",
			nil,
		)
	}
	if count > 0 {
		return nil, ErrUserAlreadyExists
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 3. Create user model
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// 4. Save to DB
	if err := s.db.Create(user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrUserAlreadyExists
		}

		return nil, apperrors.WrapApplicationError(
			err,
			apperrors.ErrInternalCategory,
			"database_error",
			"failed to create user",
			nil,
		)
	}

	return user, nil
}
