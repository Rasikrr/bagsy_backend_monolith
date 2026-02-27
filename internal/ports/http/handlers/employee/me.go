package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	coreHTTP "github.com/Rasikrr/core/http"
)

// getMe handles GET /api/v1/employees/me.
//
// @Summary      Получение профиля текущего сотрудника
// @Description  Возвращает информацию о текущем авторизованном сотруднике.
// @Tags         employees
// @Produce      json
// @Success      200  {object}  getMeResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse  "employee_not_found"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees/me [get]
func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	out, err := h.employeeUC.GetProfile(ctx, orgCtx.Employee.ID)
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, toGetMeResponse(out), http.StatusOK)
}

func toGetMeResponse(out *employeeUC.ProfileOutput) getMeResponse {
	resp := getMeResponse{
		ID:             out.ID.String(),
		Phone:          out.Phone,
		FirstName:      out.FirstName,
		LastName:       out.LastName,
		AvatarURL:      out.AvatarURL,
		OrganizationID: out.OrganizationID.String(),
		Role:           string(out.Role),
		Permissions: permissionsResponse{
			CanProvideServices:        out.Permissions.CanProvideServices,
			CanManageLocationSchedule: out.Permissions.CanManageLocationSchedule,
		},
		Active:    out.Active,
		CreatedAt: out.CreatedAt,
	}

	if out.LocationID != nil {
		s := out.LocationID.String()
		resp.LocationID = &s
	}

	return resp
}
