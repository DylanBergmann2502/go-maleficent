// internal/pkg/storage/s3_test.go
package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanFilenameRemovesPathComponents(t *testing.T) {
	assert.Equal(t, "avatar.png", cleanFilename("uploads/user/avatar.png"))
	assert.Equal(t, "avatar.png", cleanFilename(`uploads\\user\\avatar.png`))
}

func TestCleanFilenameUsesFallbackForEmptyNames(t *testing.T) {
	assert.Equal(t, "file", cleanFilename(""))
	assert.Equal(t, "file", cleanFilename("/"))
	assert.Equal(t, "file", cleanFilename("   "))
}

func TestCleanFilenameTrimsWhitespace(t *testing.T) {
	assert.Equal(t, "avatar.png", cleanFilename("  avatar.png  "))
}
