// internal/app/auth/queries/user_list_query.go
package queries

import (
	"fmt"
	"strings"

	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/errors"
	"github.com/DylanBergmann2502/go-maleficent/internal/pkg/query"
	"gorm.io/gorm"
)

var userSortFields = map[string]string{
	"id":         "id",
	"email":      "email",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

// Apply applies the translated user list query to a GORM statement.
func (params UserListQuery) Apply(db *gorm.DB) *gorm.DB {
	if params.Email != "" {
		pattern, wildcard := query.WildcardPattern(params.Email)
		if wildcard {
			db = db.Where("LOWER(email) LIKE LOWER(?)", pattern)
		} else {
			db = db.Where("email = ?", params.Email)
		}
	}
	if len(params.EmailIn) > 0 {
		db = db.Where("email IN ?", params.EmailIn)
	}
	if params.CreatedAtGTE != "" {
		db = db.Where("created_at >= ?", params.CreatedAtGTE)
	}
	if params.CreatedAtLTE != "" {
		db = db.Where("created_at <= ?", params.CreatedAtLTE)
	}
	if params.UpdatedAtGTE != "" {
		db = db.Where("updated_at >= ?", params.UpdatedAtGTE)
	}
	if params.UpdatedAtLTE != "" {
		db = db.Where("updated_at <= ?", params.UpdatedAtLTE)
	}

	for _, sortField := range params.Sort {
		direction := "ASC"
		if !sortField.Ascending {
			direction = "DESC"
		}
		db = db.Order(fmt.Sprintf("%s %s", userSortFields[sortField.Name], direction))
	}
	if len(params.Sort) == 0 {
		db = db.Order("created_at DESC").Order("id DESC")
	} else if !containsSortField(params.Sort, "id") {
		db = db.Order("id DESC")
	}

	return db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
}

func containsSortField(fields []query.SortField, name string) bool {
	for _, field := range fields {
		if strings.EqualFold(field.Name, name) {
			return true
		}
	}

	return false
}

// UserListQuery is the application-level representation of a user list
// request. It contains no Huma or HTTP-specific types.
type UserListQuery struct {
	Page         int
	PageSize     int
	Sort         []query.SortField
	Email        string
	EmailIn      []string
	CreatedAtGTE string
	CreatedAtLTE string
	UpdatedAtGTE string
	UpdatedAtLTE string
}

// NewUserListQuery translates the canonical list query input into an
// application query and rejects unsupported sort fields.
func NewUserListQuery(page, pageSize int, sortValue, email, emailIn, createdAtGTE, createdAtLTE, updatedAtGTE, updatedAtLTE string) (UserListQuery, error) {
	if page == 0 {
		page = query.DefaultPage
	}
	if pageSize == 0 {
		pageSize = query.DefaultPageSize
	}

	sortFields, err := query.ParseSort(sortValue, userSortFields)
	if err != nil {
		return UserListQuery{}, errors.NewApplicationError(
			errors.ErrValidationCategory,
			"validation_failed",
			"validation failed",
			map[string][]string{"sort": {err.Error()}},
		)
	}

	return UserListQuery{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sortFields,
		Email:        strings.TrimSpace(email),
		EmailIn:      query.ParseCSV(emailIn),
		CreatedAtGTE: createdAtGTE,
		CreatedAtLTE: createdAtLTE,
		UpdatedAtGTE: updatedAtGTE,
		UpdatedAtLTE: updatedAtLTE,
	}, nil
}
