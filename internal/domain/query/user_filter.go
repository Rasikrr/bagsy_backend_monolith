package query

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type UserFilter struct {
	NetworkCode *string
	PointCode   *string
	Roles       []enum.Role
	Phones      []string
}
