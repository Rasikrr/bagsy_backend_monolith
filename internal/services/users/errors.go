// nolint: unused
package users

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	errGetUserByPhone = coreErr.NewError("failed to get user by phone", http.StatusInternalServerError)
	errCreateUser     = coreErr.NewError("failed to create user", http.StatusInternalServerError)
	errUpdateUser     = coreErr.NewError("failed to update user", http.StatusInternalServerError)
	errSetPassword    = coreErr.NewError("failed to set password", http.StatusInternalServerError)
	errActivateUser   = coreErr.NewError("failed to activate user", http.StatusInternalServerError)
)
