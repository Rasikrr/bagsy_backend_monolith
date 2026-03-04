package catalog

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	coreHTTP "github.com/Rasikrr/core/http"
)

// createEmployeeService handles POST /api/v1/employee-services.
//
// @Summary      Привязка услуги к сотруднику
// @Description  Создаёт связь сотрудник-услуга с индивидуальной ценой. Owner — любого, Manager — staff в своей локации.
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        body  body      createEmployeeServiceRequest  true  "Данные привязки"
// @Success      201   {object}  createEmployeeServiceResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse
// @Failure      404   {object}  httputil.errorResponse
// @Failure      422   {object}  httputil.errorResponse
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employee-services [post]
func (h *Handler) createEmployeeService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req createEmployeeServiceRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.CreateEmployeeServiceInput{
		EmployeeID: req.EmployeeID,
		ServiceID:  req.ServiceID,
		Price:      req.Price,
	}

	out, err := h.catalogUseCase.CreateEmployeeService(ctx, orgCtx, input)
	if err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	coreHTTP.SendData(ctx, w, createEmployeeServiceResponse{
		ID: out.ID.String(),
	}, http.StatusCreated)
}
