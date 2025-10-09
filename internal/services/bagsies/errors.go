package bagsies

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	errInvalidParams  = coreErr.NewError("invalid params", http.StatusBadRequest)
	errCheckUserExist = coreErr.NewError("failed to check user existence", http.StatusInternalServerError)
	errCreateUser     = coreErr.NewError("failed to create user", http.StatusInternalServerError)
	errCreateBagsy    = coreErr.NewError("failed to create bagsy", http.StatusInternalServerError)
	errGetBagsies     = coreErr.NewError("failed to get bagsies", http.StatusInternalServerError)
	errDeleteBagsy    = coreErr.NewError("failed to delete bagsy", http.StatusInternalServerError)
)
