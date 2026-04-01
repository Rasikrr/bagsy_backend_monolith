package location

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// updateLocation handles PUT /api/v1/locations/{id}.
//
// @Summary      Обновление локации
// @Description  Обновляет данные локации. Все поля опциональны. Только Owner.
// @Tags         locations
// @Accept       json
// @Param        id    path  string                 true  "ID локации"
// @Param        body  body  updateLocationRequest  true  "Данные для обновления"
// @Success      204
// @Failure      400  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      410  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/locations/{id} [put]
func (h *Handler) updateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	var req updateLocationRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.UpdateLocationInput{
		ID:                  locationID,
		Name:                req.Name,
		Description:         req.Description,
		Phone:               req.Phone,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		Active:              req.Active,
		ScheduleType:        req.ScheduleType,
		SlotDurationMinutes: req.SlotDurationMinutes,
	}

	if req.Address != nil {
		input.Address = &uc.UpdateLocationAddressInput{
			City:     req.Address.City,
			Street:   req.Address.Street,
			Building: req.Address.Building,
			Details:  req.Address.Details,
		}
	}

	if err = h.locationUseCase.UpdateLocation(ctx, orgCtx, input); err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
