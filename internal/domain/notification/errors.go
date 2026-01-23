package notification

import (
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

var (
	ErrNotificationNotFound = domainErr.NewNotFoundError("notification not found", nil)
)
