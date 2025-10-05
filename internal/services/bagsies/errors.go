package bagsies

import (
	"net/http"

	"github.com/Rasikrr/core/errors"
)

var (
	errInvalidParams = errors.NewError("invalid params", http.StatusBadRequest)
)
