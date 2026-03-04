package shared

import "errors"

var (
	ErrInvalidDuration = errors.New("invalid duration")

	ErrNegativeAmount = errors.New("negative amount")
	ErrInvalidMoney   = errors.New("invalid money format")

	ErrEmptySlug = errors.New("empty slug")

	ErrInvalidPhone = errors.New("phone is invalid")
)
