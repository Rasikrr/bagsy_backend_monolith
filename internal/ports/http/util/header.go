package httputil

import (
	"errors"
	"net/http"
	"strings"
)

var ErrMissingAuthHeader = errors.New("missing or malformed authorization header")

func GetAuthHeader(r *http.Request) (string, error) {
	parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrMissingAuthHeader
	}
	return parts[1], nil
}
