package bagsies

import (
	"github.com/Rasikrr/core/errors"
	"net/http"
)

var (
	errInvalidParams = errors.NewError("invalid params", http.StatusBadRequest)
)
