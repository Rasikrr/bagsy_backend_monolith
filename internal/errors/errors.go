package errors

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	ErrUserNotFound       = coreErr.NewError("user not found", http.StatusNotFound)
	ErrSessionNotFound    = coreErr.NewError("session not found", http.StatusUnauthorized)
	ErrPhoneRequired      = coreErr.NewError("phone required", http.StatusBadRequest)
)
