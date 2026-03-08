// internal/app/users/responses/user_response.go
package responses

import (
	"github.com/maleficent/go-maleficent/internal/app/auth/models"
	"github.com/maleficent/go-maleficent/internal/pkg/responses"
)

// UserResponse defines the public user output
type UserResponse struct {
	responses.BaseResponse
	Email string `json:"email"`
}

// FromModel converts a database model to a response struct
func FromModel(u *models.User) UserResponse {
	return UserResponse{
		BaseResponse: responses.BaseResponse{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email: u.Email,
	}
}
