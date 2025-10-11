package bagsies

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// Register godoc
// @Summary Отправить код подтверждения для подтверждения записи
// @Description Отправляет код подтверждения для подтверждения записи к мастеру на whatsapp
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body confirmBagsyRequest true "Данные о услуге и юзере"
// @Success 200 {object} api.EmptySuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /api/v1/bagsies [post]
func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req confirmBagsyRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	err := c.bagsyService.SendConfirmationMessage(ctx, req.Phone, req.ServiceName)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
