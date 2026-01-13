package command

import "github.com/google/uuid"

type UpdateUserCommand struct {
	Name     string
	Surname  string
	AvatarID *uuid.UUID
}
