package httputil

import (
	"net/http"
	"strings"
)

func GetAuthHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")
	return token
}
