package media

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrMediaLimitExceeded = domainErr.NewForbiddenError("media limit exceeded")
)
