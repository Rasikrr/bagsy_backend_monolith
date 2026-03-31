package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
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

	coreHTTP.SendData(ctx, w, toGetMeResponse(out, orgCtx), http.StatusOK)
}

func toGetMeResponse(out *employeeUC.ProfileOutput, orgCtx *access.OrgContext) getMeResponse {
	resp := getMeResponse{
		ID:        out.ID.String(),
		Phone:     out.Phone,
		FirstName: out.FirstName,
		LastName:  out.LastName,
		AvatarURL: out.AvatarURL,
		Role:      string(out.Role),
		Permissions: permissionsResponse{
			CanProvideServices:        out.Permissions.CanProvideServices,
			CanManageLocationSchedule: out.Permissions.CanManageLocationSchedule,
		},
		Active:       out.Active,
		Organization: toOrganizationResponse(orgCtx),
	}

	if out.LocationID != nil {
		s := out.LocationID.String()
		resp.LocationID = &s
	}

	return resp
}

func toOrganizationResponse(orgCtx *access.OrgContext) organizationResponse {
	return organizationResponse{
		ID:           orgCtx.Organization.ID.String(),
		Name:         orgCtx.Organization.Name,
		Subscription: toSubscriptionResponse(orgCtx),
	}
}

func toSubscriptionResponse(orgCtx *access.OrgContext) subscriptionResponse {
	return subscriptionResponse{
		Plan:             orgCtx.Plan.Code.String(),
		Status:           string(orgCtx.Subscription.Status),
		CurrentPeriodEnd: orgCtx.Subscription.CurrentPeriodEnd,
		Limits:           toLimitsResponse(orgCtx),
		Features:         toFeaturesResponse(orgCtx.Plan.Code),
	}
}

func toLimitsResponse(orgCtx *access.OrgContext) limitsResponse {
	return limitsResponse{
		Locations:       toLimitValue(orgCtx.Plan.Capabilities, billing.ResourceMaxLocations, orgCtx.Subscription.LocationsUsed),
		Employees:       toLimitValue(orgCtx.Plan.Capabilities, billing.ResourceMaxEmployees, orgCtx.Subscription.EmployeesUsed),
		BookingsMonthly: limitValueResponse{Used: 0, Max: nil},
	}
}

func toLimitValue(caps access.Capabilities, resource billing.Resource, used int) limitValueResponse {
	limit, ok := caps.GetLimit(resource)
	if !ok || limit.IsUnlimited() {
		return limitValueResponse{Used: used, Max: nil}
	}
	v := limit.Value()
	return limitValueResponse{Used: used, Max: &v}
}

func toFeaturesResponse(code billing.PlanCode) featuresResponse {
	return featuresResponse{
		MultiLocation:    code.IsNetwork(),
		CustomBranding:   code.IsPoint() || code.IsNetwork(),
		APIAccess:        code.IsNetwork(),
		SMSNotifications: code.IsPoint() || code.IsNetwork(),
	}
}
