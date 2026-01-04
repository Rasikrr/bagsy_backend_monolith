package bagsy

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
)

func (c *Controller) createBagsy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createBagsyRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.bagsiesService.Create(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	// TODO: response

}
