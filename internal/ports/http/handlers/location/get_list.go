package location

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// getList handles GET /api/v1/locations.
//
// @Summary      Получение списка локаций
// @Description  Возвращает список локаций организации с пагинацией и фильтрацией. Только Owner.
// @Tags         locations
// @Produce      json
// @Param        active       query     bool      false  "Фильтр по активности"
// @Param        limit        query     int       false  "Количество записей (default: 20, max: 100)"
// @Param        offset       query     int       false  "Смещение для пагинации (default: 0)"
// @Param        order_by     query     string    false  "Поле сортировки" Enums(created_at,name) default(created_at)
// @Param        sort_order   query     string    false  "Направление сортировки" Enums(asc,desc) default(desc)
// @Success      200  {object}  getListResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/locations [get]
func (h *Handler) getList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	filter, err := parseLocationFilter(r.URL.Query())
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	filter.OrganizationID = orgCtx.Organization.ID

	page, err := h.locationUseCase.GetList(ctx, orgCtx, filter)
	if err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	now := time.Now()
	end := now.AddDate(0, 0, 7)

	items := make([]locationResponse, 0, len(page.Items))
	for _, loc := range page.Items {
		slots, sErr := h.scheduleRepo.GetLocationSlots(ctx, loc.ID, now, end)
		if sErr != nil {
			httputil.SendError(ctx, w, sErr, locationErrors)
			return
		}
		items = append(items, toLocationResponse(loc, slots))
	}

	coreHTTP.SendData(ctx, w, getListResponse{Locations: items, Total: page.Total}, http.StatusOK)
}

func parseLocationFilter(q url.Values) (*domainLoc.Filter, error) {
	filter := &domainLoc.Filter{
		Limit:     defaultLimit,
		OrderBy:   domainLoc.OrderByCreatedAt,
		SortOrder: shared.SortDesc,
	}

	if activeStr := q.Get("active"); activeStr != "" {
		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			return nil, err
		}
		filter.Active = &active
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return nil, err
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		if limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := q.Get("offset"); offsetStr != "" {
		offset, err := strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			return nil, err
		}
		filter.Offset = offset
	}

	if orderBy := q.Get("order_by"); orderBy != "" {
		parsed, err := domainLoc.ParseOrderBy(orderBy)
		if err != nil {
			return nil, err
		}
		filter.OrderBy = parsed
	}

	if sortOrder := q.Get("sort_order"); sortOrder != "" {
		parsed, err := shared.ParseSortOrder(sortOrder)
		if err != nil {
			return nil, err
		}
		filter.SortOrder = parsed
	}

	return filter, nil
}

func toLocationResponse(loc *domainLoc.Location, slots []*domainSchedule.LocationScheduleSlot) locationResponse {
	resp := locationResponse{
		ID:                  loc.ID.String(),
		CategoryID:          loc.CategoryID.String(),
		Name:                loc.Name,
		Description:         loc.Description,
		Slug:                loc.Slug.String(),
		Active:              loc.Active,
		ScheduleType:        string(loc.ScheduleType),
		SlotDurationMinutes: loc.SlotDurationMinutes.Minutes(),
		Schedule:            toScheduleSlotResponses(slots),
		CreatedAt:           loc.CreatedAt,
	}

	if loc.Phone != nil {
		p := loc.Phone.String()
		resp.Phone = &p
	}

	if loc.Address != nil {
		resp.Address = &addressResponse{
			City:     loc.Address.City,
			Street:   loc.Address.Street,
			Building: loc.Address.Building,
			Details:  loc.Address.Details,
		}
	}

	if loc.Coordinates != nil {
		resp.Coordinates = &coordsResponse{
			Latitude:  loc.Coordinates.Latitude,
			Longitude: loc.Coordinates.Longitude,
		}
	}

	return resp
}

func toScheduleSlotResponses(slots []*domainSchedule.LocationScheduleSlot) []scheduleSlotResponse {
	out := make([]scheduleSlotResponse, 0, len(slots))
	for _, s := range slots {
		out = append(out, scheduleSlotResponse{
			ID:        s.ID.String(),
			Date:      s.Date.Format("2006-01-02"),
			Type:      s.Type.String(),
			StartTime: s.StartTime.Format("15:04"),
			EndTime:   s.EndTime.Format("15:04"),
		})
	}
	return out
}
