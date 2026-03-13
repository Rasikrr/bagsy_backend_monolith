package identity

import (
	"errors"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Employee Filter & Pagination
// ─────────────────────────────────────────────────────────────────

type EmployeeFilter struct {
	OrganizationID uuid.UUID
	LocationID     *uuid.UUID
	Roles          []Role
	Search         *string
	Active         *bool
	Limit          uint64
	Offset         uint64
	OrderBy        EmployeeOrderBy
	SortOrder      shared.SortOrder
}

// ─────────────────────────────────────────────────────────────────
// OrderBy enum (whitelist for safe SQL)
// ─────────────────────────────────────────────────────────────────

type EmployeeOrderBy string

const (
	OrderByCreatedAt EmployeeOrderBy = "created_at"
	OrderByFirstName EmployeeOrderBy = "first_name"
	OrderByPhone     EmployeeOrderBy = "phone"
	OrderByRole      EmployeeOrderBy = "role"
)

var validOrderBy = map[EmployeeOrderBy]bool{
	OrderByCreatedAt: true,
	OrderByFirstName: true,
	OrderByPhone:     true,
	OrderByRole:      true,
}

var ErrInvalidOrderBy = errors.New("invalid order_by value")

func ParseEmployeeOrderBy(s string) (EmployeeOrderBy, error) {
	o := EmployeeOrderBy(s)
	if !validOrderBy[o] {
		return "", ErrInvalidOrderBy
	}
	return o, nil
}

func (o EmployeeOrderBy) String() string { return string(o) }
