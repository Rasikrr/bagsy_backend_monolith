package auth

import "errors"

var (
	errSpam         = errors.New("spam, please wait")
	errUserNotFound = errors.New("user not found")
)
