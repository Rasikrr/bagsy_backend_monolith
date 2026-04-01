package catalog

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// deleteEmployeeService handles DELETE /api/v1/employee-services/{id}.
//
// @Summary      Отвязка услуги от мастера
// @Description  Деактивирует привязку услуги к сотруднику. Owner — любого, Manager — staff своей локации.
// @Tags         catalog
// @Param        id  path  string  true  "ID привязки (employee_service)"
// @Success      204
// @Failure      400  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employee-services/{id} [delete]
func (h *Handler) deleteEmployeeService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.catalogUseCase.DeleteEmployeeService(ctx, orgCtx, id); err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
