package session

import (
	"errors"
	"net/http"
	"strings"
)

var errNoAuthHeader = errors.New("no auth header")

func GetAuthHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")
	if token == "" {
		return "", errNoAuthHeader
	}
	return token, nil
}
