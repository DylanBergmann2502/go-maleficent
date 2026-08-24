// config/config_test.go
package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setRequiredDatabaseEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("DB_HOST", "postgres")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "go_maleficent")
	t.Setenv("DB_USER", "test-user")
	t.Setenv("DB_PASSWORD", "test-password")
}

func unsetEnvironment(t *testing.T, key string) {
	t.Helper()

	value, exists := os.LookupEnv(key)
	require.NoError(t, os.Unsetenv(key))
	t.Cleanup(func() {
		if exists {
			_ = os.Setenv(key, value)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func TestLoadConfigAppendsTestDatabaseSuffix(t *testing.T) {
	setRequiredDatabaseEnvironment(t)
	t.Setenv("GO_ENV", "test")

	config, err := LoadConfig()

	require.NoError(t, err)
	assert.Equal(t, "go_maleficent_test", config.Database.Database)
	assert.Equal(t, "info", config.Log.Level)
	assert.Equal(t, "json", config.Log.Format)
}

func TestLoadConfigDoesNotDuplicateTestDatabaseSuffix(t *testing.T) {
	setRequiredDatabaseEnvironment(t)
	t.Setenv("GO_ENV", "test")
	t.Setenv("DB_NAME", "go_maleficent_test")

	config, err := LoadConfig()

	require.NoError(t, err)
	assert.Equal(t, "go_maleficent_test", config.Database.Database)
}

func TestLoadConfigUsesEnvironmentForLocalOverrides(t *testing.T) {
	setRequiredDatabaseEnvironment(t)
	t.Setenv("GO_ENV", "local")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "console")
	t.Setenv("PPROF_ENDPOINTS_ENABLED", "true")

	config, err := LoadConfig()

	require.NoError(t, err)
	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, "console", config.Log.Format)
	assert.True(t, config.Debug.PprofEndpointsEnabled)
}

func TestLoadConfigParsesEnvironmentOverrides(t *testing.T) {
	setRequiredDatabaseEnvironment(t)
	t.Setenv("GO_ENV", "production")
	t.Setenv("DB_PORT", "6432")
	t.Setenv("DB_MAX_OPEN_CONNS", "40")
	t.Setenv("DB_CONN_MAX_LIFETIME", "2m")

	config, err := LoadConfig()

	require.NoError(t, err)
	assert.Equal(t, 6432, config.Database.Port)
	assert.Equal(t, 40, config.Database.MaxOpenConns)
	assert.Equal(t, 2*time.Minute, config.Database.ConnMaxLifetime)
	assert.False(t, config.Debug.PprofEndpointsEnabled)
}

func TestLoadConfigRejectsMissingRequiredEnvironment(t *testing.T) {
	setRequiredDatabaseEnvironment(t)
	t.Setenv("GO_ENV", "test")
	unsetEnvironment(t, "DB_PASSWORD")

	_, err := LoadConfig()

	require.Error(t, err)
}
