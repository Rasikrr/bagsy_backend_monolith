package user

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrInvalidPassword       = domainErr.NewInvalidInputError("invalid password", nil)
	ErrUserNotFound          = domainErr.NewNotFoundError("user(s) not found", nil)
	ErrUserActivated         = domainErr.NewConflictError("user(s) activated", nil)
	ErrUserAlreadyExists     = domainErr.NewConflictError("user(s) already exists", nil)
	ErrUserBelongsToLocation = domainErr.NewConflictError("user(s) belong to point or network", nil)
)
