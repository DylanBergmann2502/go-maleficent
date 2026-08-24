// internal/app/auth/checks/user_exists_check.go
package checks

import (
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrEmailAlreadyRegistered identifies the user-registration conflict.
var ErrEmailAlreadyRegistered = apperrors.NewApplicationError(
	apperrors.ErrConflictCategory,
	"email_already_registered",
	"user with this email already exists",
	nil,
)

// EnsureEmailAvailable performs the workflow-level uniqueness check. The
// database unique index remains authoritative for concurrent requests.
func EnsureEmailAvailable(db *gorm.DB, email string) error {
	return EnsureEmailAvailableExcept(db, email, uuid.Nil)
}

// EnsureEmailAvailableExcept checks email availability while ignoring one
// existing user, which is needed when updating that user's email.
func EnsureEmailAvailableExcept(db *gorm.DB, email string, excludedID uuid.UUID) error {
	query := db.Model(&models.User{}).Where("email = ?", email)
	if excludedID != uuid.Nil {
		query = query.Where("id <> ?", excludedID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return apperrors.WrapApplicationError(
			err,
			apperrors.ErrUnavailableCategory,
			"database_unavailable",
			"failed to check whether the email is available",
			nil,
		)
	}

	if count > 0 {
		return ErrEmailAlreadyRegistered
	}

	return nil
}
