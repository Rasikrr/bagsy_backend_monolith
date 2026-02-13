package shared

import "errors"

var (
	ErrInvalidDuration = errors.New("invalid duration")

	ErrNegativeAmount = errors.New("negative amount")

	ErrEmptySlug = errors.New("empty slug")

	ErrInvalidPhone = errors.New("phone is invalid")
)
