package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
)

func (c *Controller) sendSmsCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req sendCodeRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}
	err := c.authService.SendCode(ctx, req.Phone)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
