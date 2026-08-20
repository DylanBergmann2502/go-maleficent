// internal/app/auth/models/user.go
package models

import (
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	models.BaseModel
	Email        string `gorm:"uniqueIndex;not null" validate:"required,email" json:"email"`
	PasswordHash string `gorm:"not null" validate:"required" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeSave enforces invariants for every workflow that persists a user.
func (u *User) BeforeSave(_ *gorm.DB) error {
	err := validator.New().Struct(u)
	if err == nil {
		return nil
	}

	return apperrors.NewApplicationError(
		apperrors.ErrValidationCategory,
		"model_validation_failed",
		"user model validation failed",
		err.Error(),
	)
}
