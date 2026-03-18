package employee

import (
	"errors"
	"net/http"

	catalogDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errEmployeeIDRequired = errors.New("employee id is required")

// getEmployeeServices handles GET /api/v1/employees/{id}/services.
//
// @Summary      Получение услуг сотрудника
// @Description  Возвращает список услуг сотрудника с его индивидуальными ценами.
// @Tags         employees
// @Produce      json
// @Param        id  path     string  true  "ID сотрудника"
// @Success      200  {object}  getEmployeeServicesResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Router       /api/v1/employees/{id}/services [get]
func (h *Handler) getEmployeeServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employeeIDStr := chi.URLParam(r, "id")
	if employeeIDStr == "" {
		httputil.SendBadRequest(ctx, w, errEmployeeIDRequired)
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	services, err := h.catalogUseCase.GetServicesByEmployee(ctx, employeeID)
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	resp := make([]employeeServiceItemResponse, 0, len(services))
	for _, s := range services {
		resp = append(resp, toEmployeeServiceItemResponse(s))
	}

	coreHTTP.SendData(ctx, w, getEmployeeServicesResponse{Services: resp}, http.StatusOK)
}

func toEmployeeServiceItemResponse(s *catalogDomain.Service) employeeServiceItemResponse {
	r := employeeServiceItemResponse{
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
		r.Price = &v
	}
	return r
}
