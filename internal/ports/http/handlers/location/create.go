package location

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	coreHTTP "github.com/Rasikrr/core/http"
)

// create handles POST /api/v1/locations.
//
// @Summary      Создание локации
// @Description  Создаёт новую точку обслуживания для организации. Проверяет лимиты тарифного плана.
// @Tags         location
// @Accept       json
// @Produce      json
// @Param        body  body      createRequest  true  "Данные локации"
// @Success      201   {object}  createResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse  "Лимит превышен или подписка приостановлена"
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/locations [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req createRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.CreateLocationInput{
		CategoryID:          req.CategoryID,
		Name:                req.Name,
		Description:         req.Description,
		Phone:               req.Phone,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		ScheduleType:        req.ScheduleType,
		SlotDurationMinutes: req.SlotDurationMinutes,
	}

	if req.Address != nil {
		input.Address = &uc.CreateLocationAddressInput{
			City:     req.Address.City,
			Street:   req.Address.Street,
			Building: req.Address.Building,
			Details:  req.Address.Details,
		}
	}

	out, err := h.locationUseCase.Create(ctx, orgCtx, input)
	if err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	coreHTTP.SendData(ctx, w, createResponse{
		ID:               out.ID.String(),
		PromptOrgProfile: out.PromptOrgProfile,
	}, http.StatusCreated)
}
