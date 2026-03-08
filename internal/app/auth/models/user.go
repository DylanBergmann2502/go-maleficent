// internal/app/users/models/user.go
package models

import (
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/models"
)

// User represents a user in the system
type User struct {
	models.BaseModel
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
