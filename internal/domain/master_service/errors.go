package masterservice

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrMasterServiceNotFound      = domainErr.NewNotFoundError("master service(s) not found", nil)
	ErrMasterServiceAlreadyExists = domainErr.NewConflictError("master service already exists", nil)
)
