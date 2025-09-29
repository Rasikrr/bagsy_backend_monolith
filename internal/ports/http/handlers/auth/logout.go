package auth

import (
	"net/http"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/cookies"
	"github.com/Rasikrr/core/api"
)

func (c *Controller) logout(w http.ResponseWriter, _ *http.Request) {
	cookies.ClearAuthTokens(w)
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
