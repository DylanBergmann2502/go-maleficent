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
	userErrors := []int{http.StatusBadRequest, http.StatusConflict, http.StatusNotFound, http.StatusUnprocessableEntity, http.StatusServiceUnavailable}

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

	huma.Register(api, huma.Operation{
		OperationID:   "list-users",
		Method:        http.MethodGet,
		Path:          "/api/v1/users",
		Summary:       "List users",
		Description:   "List all users.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusOK,
		Errors:        []int{http.StatusServiceUnavailable},
	}, userHandler.List)

	huma.Register(api, huma.Operation{
		OperationID:   "get-user",
		Method:        http.MethodGet,
		Path:          "/api/v1/users/{id}",
		Summary:       "Get a user",
		Description:   "Get a user by ID.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusOK,
		Errors:        userErrors,
	}, userHandler.Get)

	huma.Register(api, huma.Operation{
		OperationID:   "update-user",
		Method:        http.MethodPatch,
		Path:          "/api/v1/users/{id}",
		Summary:       "Update a user",
		Description:   "Update a user's email or password.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusOK,
		Errors:        userErrors,
	}, userHandler.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/api/v1/users/{id}",
		Summary:       "Delete a user",
		Description:   "Delete a user by ID.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusNotFound, http.StatusServiceUnavailable},
	}, userHandler.Delete)
}
