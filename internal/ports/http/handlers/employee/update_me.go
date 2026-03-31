package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	coreHTTP "github.com/Rasikrr/core/http"
)

// updateMe handles PUT /api/v1/employees/me.
//
// @Summary      Обновление своего профиля
// @Description  Позволяет сотруднику обновить своё имя, фамилию и аватар.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        request  body      updateMeRequest  true  "Данные для обновления"
// @Success      200  {object}  getMeResponse
// @Failure      400  {object}  httputil.errorResponse  "Неверные параметры запроса"
// @Failure      401  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse  "employee_not_found"
// @Failure      410  {object}  httputil.errorResponse  "employee_deleted"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees/me [put]
func (h *Handler) updateMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req updateMeRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.employeeUC.UpdateProfile(ctx, orgCtx.Employee.ID, employeeUC.UpdateProfileInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarID:  req.AvatarID,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, toGetMeResponse(out, orgCtx), http.StatusOK)
}
