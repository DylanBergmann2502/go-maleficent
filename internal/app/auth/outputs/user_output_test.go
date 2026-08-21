// internal/app/auth/outputs/user_output_test.go
package outputs

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DylanBergmann2502/go-maleficent/internal/app/auth/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromModelMapsPublicUserFields(t *testing.T) {
	id := uuid.New()
	createdAt := time.Date(2026, time.August, 21, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	user := &models.User{
		Email:        "user@example.com",
		PasswordHash: "secret-hash",
	}
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	output := FromModel(user)

	assert.Equal(t, id, output.ID)
	assert.Equal(t, createdAt, output.CreatedAt)
	assert.Equal(t, updatedAt, output.UpdatedAt)
	assert.Equal(t, "user@example.com", output.Email)

	encoded, err := json.Marshal(output)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "password")
}
