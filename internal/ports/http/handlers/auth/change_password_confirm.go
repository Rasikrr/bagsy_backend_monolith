package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

func (c *Controller) changePasswordConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req passwordChangeConfirmRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	err := c.authService.ChangePassword(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("password changed"), http.StatusOK)
}
