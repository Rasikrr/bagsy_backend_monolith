package jwt

import "errors"

var (
	errInvalidToken      = errors.New("invalid token")
	errTokenNotValid     = errors.New("token is not valid")
	errUnexpectedSigning = errors.New("unexpected signing method")
)
