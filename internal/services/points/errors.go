package points

import "errors"

var (
	errPointNotFound      = errors.New("point not found")
	errPointAlreadyExists = errors.New("point already exists")
)
