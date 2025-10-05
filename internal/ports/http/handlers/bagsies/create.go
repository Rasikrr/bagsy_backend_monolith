package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
)

// Register godoc
// @Summary Создание записи
// @Description Создает запись к мастеру с определенным кодом точки
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body createBagsyRequest true "Данные для создания записи"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/bagsies/create [post]
func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, err)
		return
	}

	if by.PointCode == "" {
		api.SendError(w, err)
		return
	}

	var req createBagsyRequest
	if err = api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}

	params := req.toParams()

	err = c.bagsyService.Create(ctx, params)
	if err != nil {
		return
	}

	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
