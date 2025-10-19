package httputil

import (
	"net/http"
	"strings"

	"github.com/Rasikrr/core/errors"
)

var errNoAuthHeader = errors.NewError("no auth header", http.StatusUnauthorized)

func GetAuthHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")
	if token == "" {
		return "", errNoAuthHeader
	}
	return token, nil
}
