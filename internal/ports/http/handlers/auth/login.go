package auth

import (
	"github.com/Rasikrr/core/api"
	"net/http"
)

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := api.GetData(r, req); err != nil {
		api.SendError(w, err)
		return
	}

	tokens, err := c.authService.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
	}

	resp := loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	api.SendData(w, resp, http.StatusOK)
}
