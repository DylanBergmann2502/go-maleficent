// internal/app/users/services/user_service.go
package services

import (
	"errors"

	"github.com/maleficent/go-maleficent/internal/app/auth/models"
	"github.com/maleficent/go-maleficent/internal/app/auth/utils"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
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
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
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
		return nil, err
	}

	return user, nil
}
