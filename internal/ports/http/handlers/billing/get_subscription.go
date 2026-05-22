package billing

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainBilling "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/billing"
	coreHTTP "github.com/Rasikrr/core/http"
)

// getSubscription handles GET /api/v1/subscription.
//
// @Summary      Получение текущей подписки
// @Description  Возвращает подписку текущей организации с информацией о плане.
// @Tags         billing
// @Produce      json
// @Success      200  {object}  getSubscriptionResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/subscription [get]
func (h *Handler) getSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	out, err := h.billingUC.GetSubscription(ctx, orgCtx)
	if err != nil {
		httputil.SendError(ctx, w, err, billingErrors)
		return
	}

	coreHTTP.SendData(ctx, w, getSubscriptionResponse{
		Subscription: toSubscriptionResponse(out),
	}, http.StatusOK)
}

func toSubscriptionResponse(out *uc.SubscriptionOutput) subscriptionResponse {
	sub := out.Subscription

	resp := subscriptionResponse{
		ID:                sub.ID.String(),
		OrganizationID:    sub.OrganizationID.String(),
		Status:            sub.Status.String(),
		BillingCycle:      string(sub.BillingCycle),
		RecurringAmount:   sub.RecurringAmount.String(),
		CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
		RetryCount:        sub.RetryCount,
		CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
		Plan:              toPlanResponse(out.Plan),
	}

	if sub.CurrentPeriodStart != nil {
		v := sub.CurrentPeriodStart.Format(time.RFC3339)
		resp.CurrentPeriodStart = &v
	}
	if sub.CurrentPeriodEnd != nil {
		v := sub.CurrentPeriodEnd.Format(time.RFC3339)
		resp.CurrentPeriodEnd = &v
	}
	if sub.NextBillingAt != nil {
		v := sub.NextBillingAt.Format(time.RFC3339)
		resp.NextBillingAt = &v
	}
	if sub.SuspendedAt != nil {
		v := sub.SuspendedAt.Format(time.RFC3339)
		resp.SuspendedAt = &v
	}
	if sub.CanceledAt != nil {
		v := sub.CanceledAt.Format(time.RFC3339)
		resp.CanceledAt = &v
	}
	if sub.DataDeleteAt != nil {
		v := sub.DataDeleteAt.Format(time.RFC3339)
		resp.DataDeleteAt = &v
	}

	return resp
}

func toPlanResponse(plan *domainBilling.Plan) planResponse {
	caps := make([]capabilityResponse, 0, len(plan.Capabilities))
	for _, c := range plan.Capabilities {
		cr := capabilityResponse{
			Resource: string(c.Resource),
		}
		if !c.Limit.IsUnlimited() {
			v := c.Limit.Value()
			cr.Limit = &v
		}
		caps = append(caps, cr)
	}

	return planResponse{
		ID:           plan.ID.String(),
		Code:         plan.Code.String(),
		Name:         plan.Name,
		Description:  plan.Description,
		PriceMonthly: plan.PriceMonthly.String(),
		PriceAnnual:  plan.PriceAnnual.String(),
		Capabilities: caps,
	}
}
