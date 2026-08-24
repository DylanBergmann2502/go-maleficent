// internal/app/auth/services/user_service.go
package services

import (
	stderrors "errors"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/checks"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/forms"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/utils"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/jobs/tasks"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db        *gorm.DB
	validator *validator.Validate
	jobs      *jobs.Client
}

func NewUserService(db *gorm.DB, v *validator.Validate, jobClient *jobs.Client) *UserService {
	return &UserService{db: db, validator: v, jobs: jobClient}
}

// CreateUser handles the business logic of registering a new user
func (s *UserService) CreateUser(email, password string) (*models.User, error) {
	form := forms.NewCreateUserForm(email, password)
	if err := form.Validate(s.validator); err != nil {
		return nil, err
	}
	if err := checks.EnsureEmailAvailable(s.db, form.Email); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(form.Password)
	if err != nil {
		return nil, apperrors.WrapApplicationError(
			err,
			apperrors.ErrInternalCategory,
			"password_hash_failed",
			"failed to hash password",
			nil,
		)
	}

	user := &models.User{
		Email:        form.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.db.Create(user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, checks.ErrEmailAlreadyRegistered
		}

		var appErr *apperrors.ApplicationError
		if stderrors.As(err, &appErr) {
			return nil, err
		}

		return nil, apperrors.WrapApplicationError(
			err,
			apperrors.ErrInternalCategory,
			"database_error",
			"failed to create user",
			nil,
		)
	}

	if s.jobs != nil {
		if _, err := s.jobs.Enqueue(tasks.NewExampleTask("user created")); err != nil {
			return nil, apperrors.WrapApplicationError(
				err,
				apperrors.ErrInternalCategory,
				"background_job_enqueue_failed",
				"failed to enqueue user-created task",
				nil,
			)
		}
	}

	return user, nil
}

// ListUsers returns all users without exposing password hashes.
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Order("created_at ASC").Find(&users).Error; err != nil {
		return nil, apperrors.WrapApplicationError(err, apperrors.ErrUnavailableCategory, "database_unavailable", "failed to list users", nil)
	}

	return users, nil
}

// GetUser returns a user by resource ID.
func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userNotFoundError()
		}

		return nil, apperrors.WrapApplicationError(err, apperrors.ErrUnavailableCategory, "database_unavailable", "failed to get user", nil)
	}

	return &user, nil
}

// UpdateUser applies a partial update to a user.
func (s *UserService) UpdateUser(id uuid.UUID, email, password *string) (*models.User, error) {
	form := forms.NewUpdateUserForm(email, password)
	if err := form.Validate(s.validator); err != nil {
		return nil, err
	}

	user, err := s.GetUser(id)
	if err != nil {
		return nil, err
	}

	if form.Email != nil && *form.Email != user.Email {
		if err := checks.EnsureEmailAvailableExcept(s.db, *form.Email, id); err != nil {
			return nil, err
		}
		user.Email = *form.Email
	}

	if form.Password != nil {
		hashedPassword, err := utils.HashPassword(*form.Password)
		if err != nil {
			return nil, apperrors.WrapApplicationError(err, apperrors.ErrInternalCategory, "password_hash_failed", "failed to hash password", nil)
		}
		user.PasswordHash = hashedPassword
	}

	if err := s.db.Save(user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, checks.ErrEmailAlreadyRegistered
		}
		var appErr *apperrors.ApplicationError
		if stderrors.As(err, &appErr) {
			return nil, err
		}
		return nil, apperrors.WrapApplicationError(err, apperrors.ErrInternalCategory, "database_error", "failed to update user", nil)
	}

	return user, nil
}

// DeleteUser removes a user by resource ID.
func (s *UserService) DeleteUser(id uuid.UUID) error {
	result := s.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return apperrors.WrapApplicationError(result.Error, apperrors.ErrInternalCategory, "database_error", "failed to delete user", nil)
	}
	if result.RowsAffected == 0 {
		return userNotFoundError()
	}

	return nil
}

func userNotFoundError() error {
	return apperrors.NewApplicationError(apperrors.ErrNotFoundCategory, "user_not_found", "user not found", nil)
}
