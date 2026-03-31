package location

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	coreHTTP "github.com/Rasikrr/core/http"
)

// updateOrganization handles PUT /api/v1/organizations/me.
//
// @Summary      Обновление профиля организации
// @Description  Позволяет владельцу обновить имя и описание организации (сети).
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        request  body      updateOrganizationRequest  true  "Данные для обновления"
// @Success      200  {object}  updateOrganizationResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse  "permission_denied"
// @Failure      404  {object}  httputil.errorResponse  "organization_not_found"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/organizations/me [put]
func (h *Handler) updateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req updateOrganizationRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err := h.locationUseCase.UpdateOrganization(ctx, orgCtx, uc.UpdateOrganizationInput{
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	coreHTTP.SendData(ctx, w, updateOrganizationResponse{
		ID:          orgCtx.Organization.ID.String(),
		Name:        req.Name,
		Description: req.Description,
	}, http.StatusOK)
}
