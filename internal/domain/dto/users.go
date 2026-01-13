package dto

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"

type PaginatedUsers struct {
	Users []*UserWithAvatar
	Total int
}

type UserWithAvatar struct {
	*entity.User

	AvatarURL *string
}
