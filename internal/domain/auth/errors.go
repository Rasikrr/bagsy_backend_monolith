package auth

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

// JWT and Authentication domain errors
var (
	// JWT token errors
	ErrInvalidToken            = errors.NewUnauthorizedError("invalid token")
	ErrTokenExpired            = errors.NewUnauthorizedError("token has expired")
	ErrTokenNotValid           = errors.NewUnauthorizedError("token is not valid")
	ErrUnexpectedSigningMethod = errors.NewUnauthorizedError("unexpected signing method")
	ErrRefreshTokenRequired    = errors.NewInvalidInputError("refresh token required", nil)
	ErrAccessTokenRequired     = errors.NewInvalidInputError("access token required", nil)
	ErrTokenRevoked            = errors.NewUnauthorizedError("token has been revoked")
	ErrInvalidTokenClaims      = errors.NewUnauthorizedError("invalid token claims")

	// Authentication errors
	ErrInvalidCredentials = errors.NewUnauthorizedError("invalid credentials")
	ErrInvalidPassword    = errors.NewUnauthorizedError("invalid password")
	ErrUserInactive       = errors.NewForbiddenError("user is inactive")
	ErrNoPassword         = errors.NewForbiddenError("user has no password set")

	ErrUnknownRole = errors.NewUnauthorizedError("unknown role")

	// Verification code errors
	ErrTooManyVerificationAttempts = errors.NewForbiddenError("maximum attempts reached. Please restart your registration")
	ErrInvalidVerificationCode     = errors.NewInvalidInputError("invalid verification code", nil)
)
