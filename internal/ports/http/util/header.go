package util

import (
	"net/http"
	"strings"
)

func GetAuthHeader(r *http.Request) (string, error) {
	bearerToken := strings.Split(r.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 {

	}
	return bearerToken[1], nil
}
