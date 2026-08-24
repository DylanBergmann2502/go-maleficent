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
