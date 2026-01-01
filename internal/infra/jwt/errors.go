package jwt

import "github.com/cockroachdb/errors"

// JWT ошибки - простые errors без зависимостей от domain
// Service слой должен маппить эти ошибки в доменные

var (
	// ErrTokenExpired токен истек
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken невалидный токен
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSignature неверная подпись
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrMalformedToken некорректный формат токена
	ErrMalformedToken = errors.New("malformed token")

	// ErrTokenNotValid токен не валиден
	ErrTokenNotValid = errors.New("token not valid")

	// ErrUnexpectedSigningMethod неожиданный метод подписи
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrInvalidTokenClaims невалидные claims токена
	ErrInvalidTokenClaims = errors.New("invalid token claims")

	// ErrRegistrationTokenExpired registration токен истек
	ErrRegistrationTokenExpired = errors.New("registration token expired")

	// ErrInvalidRegistrationToken невалидный registration токен
	ErrInvalidRegistrationToken = errors.New("invalid registration token")

	// ErrMissingRequiredClaim отсутствует обязательный claim
	ErrMissingRequiredClaim = errors.New("missing required claim")
)
