// nolint: unused
package auth

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	errGetUser                 = coreErr.NewError("failed to get user", http.StatusInternalServerError)
	errInvalidToken            = coreErr.NewError("invalid token", http.StatusUnauthorized)
	errSpam                    = coreErr.NewError("spam detected", http.StatusTooManyRequests)
	errNoAccess                = coreErr.NewError("no access", http.StatusForbidden)
	errInvalidPassword         = coreErr.NewError("invalid password", http.StatusUnauthorized)
	errAccessTokenNotAllowed   = coreErr.NewError("access token is not allowed", http.StatusBadRequest)
	errRefreshTokenNotAllowed  = coreErr.NewError("refresh token is not allowed", http.StatusBadRequest)
	errGenerateTokens          = coreErr.NewError("failed to generate tokens", http.StatusInternalServerError)
	errHashingFailed           = coreErr.NewError("hashing failed", http.StatusInternalServerError)
	errSetPassword             = coreErr.NewError("failed to set password", http.StatusInternalServerError)
	errActivateUser            = coreErr.NewError("failed to activate user", http.StatusInternalServerError)
	errRegistrationConfirm     = coreErr.NewError("registration confirm failed", http.StatusInternalServerError)
	errGenerateRegistrationURL = coreErr.NewError("failed to generate registration url", http.StatusInternalServerError)
)
