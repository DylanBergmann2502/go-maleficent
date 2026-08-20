// internal/app/auth/inputs/user_input.go
package inputs

// CreateUserInput is the Huma/OpenAPI input for user registration.
// Workflow-specific validation remains in the handler/service layer.
type CreateUserInput struct {
	Body struct {
		Email    string `json:"email" format:"email" minLength:"1" doc:"User email address" example:"user@example.com"`
		Password string `json:"password" minLength:"8" maxLength:"72" doc:"User password" example:"correct horse battery staple"`
	} `json:"body"`
}
