package media

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
)

func (c *Controller) getUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req uploadURLRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

}
