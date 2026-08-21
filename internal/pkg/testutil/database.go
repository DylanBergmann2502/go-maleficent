// internal/pkg/testutil/database.go
package testutil

import (
	"os"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/configs"
	databasepkg "github.com/DylanBergmann2502/go-maleficent/internal/pkg/database"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewDatabaseClient opens the shared PostgreSQL test database client.
func NewDatabaseClient(t *testing.T) *databasepkg.Database {
	t.Helper()

	require.Equal(t, "test", os.Getenv("GO_ENV"), "tests must run with GO_ENV=test")

	config, err := configs.LoadConfig()
	require.NoError(t, err)

	database, err := databasepkg.NewDatabase(&config.Database, zap.NewNop())
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, database.Close())
	})

	return database
}

// NewDatabase opens the shared PostgreSQL test database.
func NewDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	return NewDatabaseClient(t).DB
}

// Begin starts a transaction that is rolled back when the test finishes.
// Tests can seed records through the returned handle without leaving data in
// the shared test database.
func Begin(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	transaction := db.Begin()
	require.NoError(t, transaction.Error)

	t.Cleanup(func() {
		require.NoError(t, transaction.Rollback().Error)
	})

	return transaction
}
