// nolint: unused
package auth

import "errors"

var (
	errSpam            = errors.New("spam, please wait")
	errNoAccess        = errors.New("no access")
	errInvalidPassword = errors.New("invalid password")
)
