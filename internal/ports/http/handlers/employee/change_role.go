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

// nolint: dupl
// changeRole handles PATCH /api/v1/employees/{id}/role.
//
// @Summary      Смена роли сотрудника
// @Description  Только owner может менять роли. Нельзя менять свою роль и назначать роль owner.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id       path      string            true  "UUID сотрудника"
// @Param        request  body      changeRoleRequest true  "Новая роль"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.errorResponse  "invalid_role / cannot_set_owner_role"
// @Failure      403  {object}  httputil.errorResponse  "permission_denied / cannot_modify_self"
// @Failure      404  {object}  httputil.errorResponse  "employee_not_found"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees/{id}/role [patch]
func (h *Handler) changeRole(w http.ResponseWriter, r *http.Request) {
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

	var req changeRoleRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	err = h.employeeUC.ChangeRole(ctx, orgCtx, id, employeeUC.ChangeRoleInput{
		Role: req.Role,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, map[string]string{"message": "role_changed"}, http.StatusOK)
}
