// internal/pkg/query/query.go
package query

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

var sortPattern = regexp.MustCompile(`^-?[A-Za-z][A-Za-z0-9_]*(,-?[A-Za-z][A-Za-z0-9_]*)*$`)

// SortField represents one field in the canonical list sort format.
type SortField struct {
	Name      string
	Ascending bool
}

// Pagination contains page-based pagination metadata.
type Pagination struct {
	Page       int
	PageSize   int
	TotalCount int64
	TotalPages int
	HasNext    bool
	HasPrev    bool
	NextPage   *int
	PrevPage   *int
}

// NewPagination builds pagination metadata using one-indexed pages.
func NewPagination(page, pageSize int, totalCount int64) Pagination {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	result := Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}

	if result.HasPrev {
		previous := page - 1
		result.PrevPage = &previous
	}
	if result.HasNext {
		next := page + 1
		result.NextPage = &next
	}

	return result
}

// ParseSort parses the comma-separated sort format and validates
// every field against the supplied public-to-database field mapping.
func ParseSort(value string, allowed map[string]string) ([]SortField, error) {
	if value == "" {
		return nil, nil
	}
	if !sortPattern.MatchString(value) {
		return nil, fmt.Errorf("invalid sort format")
	}

	fields := make([]SortField, 0, strings.Count(value, ",")+1)
	for _, item := range strings.Split(value, ",") {
		ascending := true
		name := item
		if strings.HasPrefix(name, "-") {
			ascending = false
			name = strings.TrimPrefix(name, "-")
		}
		if _, ok := allowed[name]; !ok {
			return nil, fmt.Errorf("sort field %q is not supported", name)
		}
		fields = append(fields, SortField{Name: name, Ascending: ascending})
	}

	return fields, nil
}

// ParseCSV parses the canonical comma-separated IN filter format.
func ParseCSV(value string) []string {
	values := strings.Split(value, ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}

	return result
}

// WildcardPattern converts the '*' wildcard syntax to SQL LIKE
// syntax. It returns false when the value is an exact-match value.
func WildcardPattern(value string) (string, bool) {
	if !strings.Contains(value, "*") {
		return value, false
	}

	return strings.ReplaceAll(value, "*", "%"), true
}
