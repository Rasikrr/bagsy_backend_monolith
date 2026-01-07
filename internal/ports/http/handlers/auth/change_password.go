package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

func (c *Controller) changePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req changePasswordRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.authService.SendPasswordChangeLink(ctx, req.Phone)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("link sent"), http.StatusOK)
}
