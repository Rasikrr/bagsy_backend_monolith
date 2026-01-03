package query

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type UserFilter struct {
	NetworkCode *string
	PointCode   *string
	Roles       []enum.Role
	Phones      []string
	Limit       uint64
	Offset      uint64
	SortBy      string
	SortOrder   enum.SortOrder
}
