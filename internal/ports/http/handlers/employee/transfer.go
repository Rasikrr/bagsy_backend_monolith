// nolint: dupl
package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// transferEmployee handles POST /api/v1/employees/{id}/transfer.
//
// @Summary      Перевод сотрудника в другую локацию
// @Description  Owner может перевести любого сотрудника. Manager — только staff своей локации.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id       path      string            true  "UUID сотрудника"
// @Param        request  body      transferRequest   true  "Данные для перевода"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse  "permission_denied / cannot_modify_self"
// @Failure      404  {object}  httputil.errorResponse  "employee_not_found"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees/{id}/transfer [post]
func (h *Handler) transferEmployee(w http.ResponseWriter, r *http.Request) {
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

	var req transferRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	err = h.employeeUC.TransferEmployee(ctx, orgCtx, id, employeeUC.TransferInput{
		LocationID: req.LocationID,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, map[string]string{"message": "employee_transferred"}, http.StatusOK)
}
