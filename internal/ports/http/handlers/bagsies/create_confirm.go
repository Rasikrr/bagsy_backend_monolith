package bagsies

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// Register godoc
// @Summary Создание записи
// @Description Создает запись к мастеру с определенным кодом точки
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body createBagsyRequest true "Данные для создания записи"
// @Success 200 {object} api.EmptySuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /api/v1/bagsies/confirm [post]
func (c *Controller) createConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createBagsyRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}

	err := c.bagsyService.Create(ctx, req.toParams())
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.SendData(w,
		api.NewEmptySuccessResponse("bagsy created"),
		http.StatusCreated,
	)
}
