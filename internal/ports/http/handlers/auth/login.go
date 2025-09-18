package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
)

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := api.GetData(r, req); err != nil {
		api.SendError(w, err)
		return
	}

	if err := req.validate(); err != nil {
		api.SendError(w, err)
		return
	}

	tokens, err := c.authService.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
		return
	}

	resp := loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	api.SendData(w, resp, http.StatusOK)
}
