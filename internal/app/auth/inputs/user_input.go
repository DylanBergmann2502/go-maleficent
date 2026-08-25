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

// ListUsersInput contains the list endpoint's transport parameters.
type ListUsersInput struct {
	Page         int    `query:"page" minimum:"1" default:"1" doc:"Page number (1-indexed)"`
	PageSize     int    `query:"page_size" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
	Sort         string `query:"sort" pattern:"^-?[A-Za-z][A-Za-z0-9_]*(,-?[A-Za-z][A-Za-z0-9_]*)*$" doc:"Comma-separated sort fields; prefix a field with - for descending order" example:"-created_at,email"`
	Email        string `query:"email" doc:"Filter by exact email or wildcard pattern" example:"*@example.com"`
	EmailIn      string `query:"email__in" doc:"Filter by comma-separated exact email values" example:"one@example.com,two@example.com"`
	CreatedAtGTE string `query:"created_at__gte" format:"date-time" doc:"Filter users created at or after this timestamp"`
	CreatedAtLTE string `query:"created_at__lte" format:"date-time" doc:"Filter users created at or before this timestamp"`
	UpdatedAtGTE string `query:"updated_at__gte" format:"date-time" doc:"Filter users updated at or after this timestamp"`
	UpdatedAtLTE string `query:"updated_at__lte" format:"date-time" doc:"Filter users updated at or before this timestamp"`
}

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
