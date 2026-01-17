package actor

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrUnknownRole = domainErr.NewUnauthorizedError("unknown role")
)
