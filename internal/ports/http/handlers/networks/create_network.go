package networks

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// Пока не юзаем, флоу может поменяться
func (c *Controller) createNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createNetworkRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.networksService.Create(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("network created"), http.StatusCreated)
}
