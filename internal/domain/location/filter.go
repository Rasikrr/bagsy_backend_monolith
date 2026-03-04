package location

import (
	"errors"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Location Filter & Pagination
// ─────────────────────────────────────────────────────────────────

type Filter struct {
	OrganizationID uuid.UUID
	Active         *bool
	Limit          uint64
	Offset         uint64
	OrderBy        OrderBy
	SortOrder      shared.SortOrder
}

// ─────────────────────────────────────────────────────────────────
// OrderBy enum (whitelist for safe SQL)
// ─────────────────────────────────────────────────────────────────

type OrderBy string

const (
	OrderByCreatedAt OrderBy = "created_at"
	OrderByName      OrderBy = "name"
)

var validOrderBy = map[OrderBy]bool{
	OrderByCreatedAt: true,
	OrderByName:      true,
}

var ErrInvalidOrderBy = errors.New("invalid order_by value")

func ParseOrderBy(s string) (OrderBy, error) {
	o := OrderBy(s)
	if !validOrderBy[o] {
		return "", ErrInvalidOrderBy
	}
	return o, nil
}

func (o OrderBy) String() string { return string(o) }
