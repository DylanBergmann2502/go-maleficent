// internal/pkg/query/query_test.go
package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSort(t *testing.T) {
	fields, err := ParseSort("-created_at,email", map[string]string{
		"created_at": "created_at",
		"email":      "email",
	})

	require.NoError(t, err)
	assert.Equal(t, []SortField{
		{Name: "created_at", Ascending: false},
		{Name: "email", Ascending: true},
	}, fields)
}

func TestParseSortRejectsUnknownFields(t *testing.T) {
	_, err := ParseSort("password_hash", map[string]string{"email": "email"})

	assert.EqualError(t, err, `sort field "password_hash" is not supported`)
}

func TestNewPagination(t *testing.T) {
	pagination := NewPagination(2, 20, 45)

	assert.Equal(t, 3, pagination.TotalPages)
	assert.True(t, pagination.HasPrev)
	assert.True(t, pagination.HasNext)
	require.NotNil(t, pagination.PrevPage)
	require.NotNil(t, pagination.NextPage)
	assert.Equal(t, 1, *pagination.PrevPage)
	assert.Equal(t, 3, *pagination.NextPage)
}

func TestWildcardPattern(t *testing.T) {
	pattern, wildcard := WildcardPattern("*Example*")

	assert.True(t, wildcard)
	assert.Equal(t, "%Example%", pattern)
}
