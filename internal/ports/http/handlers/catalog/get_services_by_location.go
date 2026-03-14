package catalog

import (
	"errors"
	"net/http"

	catalogDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errLocationIDRequired = errors.New("location_id is required")

// getServicesByLocation handles GET /api/v1/services?location_id=.
//
// @Summary      Получение услуг по локации
// @Description  Возвращает список услуг для указанной локации.
// @Tags         catalog
// @Produce      json
// @Param        id  path     string  true  "ID локации"
// @Success      200  {object}  getServicesByLocationResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/services/{id} [get]
func (h *Handler) getServicesByLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationIDStr := chi.URLParam(r, "id")
	if locationIDStr == "" {
		httputil.SendBadRequest(ctx, w, errLocationIDRequired)
		return
	}

	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	services, err := h.catalogUseCase.GetServicesByLocation(ctx, locationID)
	if err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	resp := make([]serviceResponse, 0, len(services))
	for _, s := range services {
		resp = append(resp, toServiceResponse(s))
	}

	coreHTTP.SendData(ctx, w, getServicesByLocationResponse{Services: resp}, http.StatusOK)
}

func toServiceResponse(s *catalogDomain.Service) serviceResponse {
	r := serviceResponse{
		ID:              s.ID.String(),
		CategoryID:      s.CategoryID.String(),
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes.Minutes(),
		Color:           string(s.Color),
		SortOrder:       s.SortOrder,
		Active:          s.Active,
	}
	if s.MinPrice != nil {
		v := s.MinPrice.Amount().IntPart()
		r.MinPrice = &v
	}
	if s.MaxPrice != nil {
		v := s.MaxPrice.Amount().IntPart()
		r.MaxPrice = &v
	}
	return r
}
