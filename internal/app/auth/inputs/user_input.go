// internal/app/auth/inputs/user_input.go
package inputs

import "github.com/google/uuid"

// CreateUserInput is the Huma/OpenAPI input for user registration.
// Workflow-specific validation remains in the handler/service layer.
type CreateUserInput struct {
	Body struct {
		Email    string `json:"email" format:"email" minLength:"1" doc:"User email address" example:"user@example.com"`
		Password string `json:"password" minLength:"8" maxLength:"72" doc:"User password" example:"correct horse battery staple"`
	} `json:"body"`
}

// ListUsersInput contains filters for listing users.
type ListUsersInput struct{}

// GetUserInput identifies a user resource.
type GetUserInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"User ID"`
}

// UpdateUserInput contains a partial user update.
type UpdateUserInput struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"User ID"`
	Body struct {
		Email    *string `json:"email,omitempty" format:"email" minLength:"1" doc:"User email address" example:"user@example.com"`
		Password *string `json:"password,omitempty" minLength:"8" maxLength:"72" doc:"Replacement user password" example:"correct horse battery staple"`
	} `json:"body"`
}

// DeleteUserInput identifies a user resource to delete.
type DeleteUserInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"User ID"`
}
