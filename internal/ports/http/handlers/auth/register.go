package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
)

func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := api.GetData(r, &req); err != nil {
		return
	}
	if err := req.validate(); err != nil {
		api.SendError(w, err)
		return
	}

	err := c.usersService.Create(r.Context(), req.convert())
	if err != nil {
		api.SendError(w, err)
		return
	}
	// TODO: send temporary link to whatsapp
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
