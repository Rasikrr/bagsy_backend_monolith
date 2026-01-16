package bagsy

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrBagsyNotFound              = domainErr.NewNotFoundError("bagsies(ies) not found", nil)
	ErrBagsyTimeIsAlreadyOccupied = domainErr.NewConflictError("time is already occupied", nil)
)
