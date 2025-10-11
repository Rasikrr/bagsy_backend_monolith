// nolint: unused
package forms

import (
	"net/http"

	coreErr "github.com/Rasikrr/core/errors"
)

var (
	errCreateClient = coreErr.NewError("failed to create client", http.StatusInternalServerError)
)
