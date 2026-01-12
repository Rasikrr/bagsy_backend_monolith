package dto

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"

type PaginatedUsers struct {
	Users []*entity.User
	Total int
}
