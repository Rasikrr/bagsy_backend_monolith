// nolint: unused
package users

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	errDeleteUnactivatedUsers = coreErr.NewError("failed to delete unactivated users", http.StatusInternalServerError)
	errValidateParams         = coreErr.NewError("failed to validate params", http.StatusBadRequest)
	errNoPermission           = coreErr.NewError("no permission to create user with higher role", http.StatusForbidden)
	errUserAlreadyExists      = coreErr.NewError("user already exists", http.StatusBadRequest)
	errGetUser                = coreErr.NewError("failed to get user", http.StatusInternalServerError)
	errCreateUser             = coreErr.NewError("failed to create user", http.StatusInternalServerError)
	errUpdateUser             = coreErr.NewError("failed to update user", http.StatusInternalServerError)
	errSetPassword            = coreErr.NewError("failed to set password", http.StatusInternalServerError)
	errActivateUser           = coreErr.NewError("failed to activate user", http.StatusInternalServerError)
)
