package errors

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	ErrSessionNotFound = coreErr.NewError("session not found", http.StatusUnauthorized)
	ErrPhoneRequired   = coreErr.NewError("phone required", http.StatusBadRequest)
)
