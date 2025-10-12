package points

import (
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"net/http"

	"github.com/Rasikrr/core/api"
	"github.com/go-chi/chi/v5"
)

func (c *Controller) getByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "code")
	point, err := c.service.GetByCode(ctx, code)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, point, http.StatusOK)
}

func (c *Controller) updateByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	var req UpdatePointRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}
	point := req.ToEntity(code)
	if err := c.service.UpdateByCode(ctx, code, point); err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, point, http.StatusOK)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sess, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, err)
		return
	}

	var req CreatePointRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}

	point := req.ToEntity(sess.Phone)
	err = c.service.Create(ctx, point)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.SendData(w, api.NewEmptySuccessResponse("created"), http.StatusCreated)
}

func (c *Controller) deleteByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "code")
	if err := c.service.DeleteByCode(ctx, code); err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, nil, http.StatusNoContent)
}
