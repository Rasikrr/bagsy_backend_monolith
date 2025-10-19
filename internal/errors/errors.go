package errors

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	ErrNoPermission        = coreErr.NewError("no permission", http.StatusForbidden)
	ErrUserNotFound        = coreErr.NewError("user not found", http.StatusNotFound)
	ErrSessionNotFound     = coreErr.NewError("session not found", http.StatusUnauthorized)
	ErrPhoneRequired       = coreErr.NewError("phone required", http.StatusBadRequest)
	ErrNetworkCodeRequired = coreErr.NewError("network code required", http.StatusBadRequest)
	ErrPointCodeRequired   = coreErr.NewError("point code required", http.StatusBadRequest)
)

var (
	ErrNotImplemented = coreErr.NewError("not implemented", http.StatusNotImplemented)
)
