// internal/app/auth/outputs/user_output.go
package outputs

import (
	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	apiresponses "github.com/DylanBergmann2502/go-maleficent/internal/pkg/api/responses"
)

// UserOutput is the public representation of a user.
type UserOutput struct {
	apiresponses.BaseResponse
	Email string `json:"email"`
}

// CreateUserOutput is the Huma/OpenAPI output for user registration.
type CreateUserOutput struct {
	Body struct {
		Data UserOutput            `json:"data"`
		Meta apiresponses.MetaData `json:"meta"`
	} `json:"body"`
}

// ListUsersOutput is the Huma/OpenAPI output for listing users.
type ListUsersOutput struct {
	Body struct {
		Data []UserOutput          `json:"data"`
		Meta apiresponses.MetaData `json:"meta"`
	} `json:"body"`
}

// GetUserOutput is the Huma/OpenAPI output for retrieving a user.
type GetUserOutput struct {
	Body struct {
		Data UserOutput            `json:"data"`
		Meta apiresponses.MetaData `json:"meta"`
	} `json:"body"`
}

// UpdateUserOutput is the Huma/OpenAPI output for updating a user.
type UpdateUserOutput struct {
	Body struct {
		Data UserOutput            `json:"data"`
		Meta apiresponses.MetaData `json:"meta"`
	} `json:"body"`
}

// DeleteUserOutput represents an empty successful delete response.
type DeleteUserOutput struct{}

// FromModel converts a database model to a public output.
func FromModel(user *models.User) UserOutput {
	return UserOutput{
		apiresponses.BaseResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		user.Email,
	}
}
