// internal/pkg/models/base_test.go
package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseModelBeforeCreateGeneratesUUIDV7(t *testing.T) {
	base := &BaseModel{}

	require.NoError(t, base.BeforeCreate(nil))

	assert.NotEqual(t, uuid.Nil, base.ID)
	assert.Equal(t, uuid.Version(7), base.ID.Version())
}

func TestBaseModelBeforeCreatePreservesExistingID(t *testing.T) {
	existingID := uuid.New()
	base := &BaseModel{ID: existingID}

	require.NoError(t, base.BeforeCreate(nil))

	assert.Equal(t, existingID, base.ID)
}
