package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// activateEmployee handles POST /api/v1/employees/{id}/activate.
//
// @Summary      Активация сотрудника
// @Description  Owner может активировать любого. Manager — только staff своей локации.
// @Tags         employees
// @Produce      json
// @Param        id  path  string  true  "UUID сотрудника"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  httputil.errorResponse  "permission_denied / cannot_modify_self"
// @Failure      404  {object}  httputil.errorResponse  "employee_not_found"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees/{id}/activate [post]
func (h *Handler) activateEmployee(w http.ResponseWriter, r *http.Request) {
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

	if err = h.employeeUC.ActivateEmployee(ctx, orgCtx, id); err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, map[string]string{"message": "employee_activated"}, http.StatusOK)
}
