// internal/app/auth/checks/user_exists_check.go
package checks

import (
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	apperrors "github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
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
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
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
