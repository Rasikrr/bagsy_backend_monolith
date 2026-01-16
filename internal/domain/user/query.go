package user

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type Filter struct {
	NetworkCode *string
	PointCode   *string
	Roles       []Role
	Phones      []string
	Limit       uint64
	Offset      uint64
	OrderBy     string
	SortOrder   enum.SortOrder
}
