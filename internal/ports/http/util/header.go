package httputil

import (
	"net/http"
	"strings"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func GetAuthHeader(r *http.Request) (string, error) {
	bearerToken := strings.Split(r.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 {
		return "", domainErr.NewUnauthorizedError("Invalid Authorization header")
	}
	return bearerToken[1], nil
}
