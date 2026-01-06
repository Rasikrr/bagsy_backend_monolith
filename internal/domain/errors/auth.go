package errors

// JWT and Authentication domain errors
var (
	// JWT token errors
	ErrInvalidToken            = NewUnauthorizedError("invalid token")
	ErrTokenExpired            = NewUnauthorizedError("token has expired")
	ErrTokenNotValid           = NewUnauthorizedError("token is not valid")
	ErrUnexpectedSigningMethod = NewUnauthorizedError("unexpected signing method")
	ErrRefreshTokenRequired    = NewInvalidInputError("refresh token required", nil)
	ErrAccessTokenRequired     = NewInvalidInputError("access token required", nil)
	ErrTokenRevoked            = NewUnauthorizedError("token has been revoked")
	ErrInvalidTokenClaims      = NewUnauthorizedError("invalid token claims")

	// Authentication errors
	ErrInvalidCredentials = NewUnauthorizedError("invalid credentials")
	ErrInvalidPassword    = NewUnauthorizedError("invalid password")
	ErrUserInactive       = NewForbiddenError("user is inactive")
	ErrNoPassword         = NewForbiddenError("user has no password set")

	// Registration token errors
	ErrInvalidRegistrationToken = NewUnauthorizedError("invalid registration token")
	ErrRegistrationTokenExpired = NewUnauthorizedError("registration token has expired")
	ErrMissingRequiredClaim     = NewInvalidInputError("missing required claim", nil)

	// Verification code errors
	ErrTooManyVerificationAttempts = NewForbiddenError("maximum attempts reached. Please restart your registration")
	ErrInvalidVerificationCode     = NewInvalidInputError("invalid verification code", nil)
)
