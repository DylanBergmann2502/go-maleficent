// internal/app/auth/routes.go
package auth

import (
	"net/http"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/handlers"
	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers the authentication API operations and their Huma
// schemas on the shared API instance.
func RegisterRoutes(api huma.API, userHandler *handlers.UserHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/v1/users",
		Summary:       "Create a user",
		Description:   "Create a new user account.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusBadRequest, http.StatusConflict, http.StatusUnprocessableEntity},
	}, userHandler.Create)
}
