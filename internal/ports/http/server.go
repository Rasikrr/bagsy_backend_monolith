package http

import "github.com/Rasikrr/core/http"

func NewServer(
	server *http.Server,
) {
	server.WithControllers()
}
