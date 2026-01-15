package network

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrNetworkNotFound      = domainErr.NewNotFoundError("network(s) not found", nil)
	ErrNetworkAlreadyExists = domainErr.NewConflictError("network(s) already exists", nil)
)
