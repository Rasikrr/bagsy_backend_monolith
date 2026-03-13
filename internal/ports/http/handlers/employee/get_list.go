package employee

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/google/uuid"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// getList handles GET /api/v1/employees.
//
// @Summary      Получение списка сотрудников
// @Description  Возвращает список сотрудников организации с пагинацией и фильтрацией.
// @Description  Owner — все сотрудники организации. Manager — только сотрудники своей локации. Staff — доступ запрещён (403).
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        location_id  query     string    false  "UUID локации для фильтрации"
// @Param        role         query     []string  false  "Фильтр по ролям (multi-value)" collectionFormat(multi) Enums(owner,manager,staff)
// @Param        search       query     string    false  "Поиск по имени или номеру телефона (ILIKE)"
// @Param        active       query     bool      false  "Фильтр по активности"
// @Param        limit        query     int       false  "Количество записей (default: 20, max: 100)"
// @Param        offset       query     int       false  "Смещение для пагинации (default: 0)"
// @Param        order_by     query     string    false  "Поле сортировки" Enums(created_at,first_name,phone,role) default(created_at)
// @Param        sort_order   query     string    false  "Направление сортировки" Enums(asc,desc) default(desc)
// @Success      200  {object}  getListResponse
// @Failure      400  {object}  httputil.errorResponse  "Неверные параметры запроса"
// @Failure      401  {object}  httputil.errorResponse  "Требуется авторизация"
// @Failure      403  {object}  httputil.errorResponse  "Недостаточно прав или подписка приостановлена"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employees [get]
func (h *Handler) getList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	filter, err := parseEmployeeFilter(r.URL.Query())
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	result, err := h.employeeUC.GetList(ctx, orgCtx, filter)
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, toGetListResponse(result), http.StatusOK)
}

func parseEmployeeFilter(q url.Values) (*identity.EmployeeFilter, error) {
	filter := &identity.EmployeeFilter{
		Limit:     defaultLimit,
		OrderBy:   identity.OrderByCreatedAt,
		SortOrder: shared.SortDesc,
	}

	if err := parseLocationID(q, filter); err != nil {
		return nil, err
	}
	if err := parseRoles(q, filter); err != nil {
		return nil, err
	}
	parseSearch(q, filter)
	if err := parseActive(q, filter); err != nil {
		return nil, err
	}
	if err := parsePagination(q, filter); err != nil {
		return nil, err
	}
	if err := parseSorting(q, filter); err != nil {
		return nil, err
	}

	return filter, nil
}

func parseLocationID(q url.Values, filter *identity.EmployeeFilter) error {
	locID := q.Get("location_id")
	if locID == "" {
		return nil
	}
	parsed, err := uuid.Parse(locID)
	if err != nil {
		return err
	}
	filter.LocationID = &parsed
	return nil
}

func parseRoles(q url.Values, filter *identity.EmployeeFilter) error {
	roles := q["role"]
	if len(roles) == 0 {
		return nil
	}
	parsed := make([]identity.Role, 0, len(roles))
	for _, rs := range roles {
		role := identity.Role(rs)
		if !role.IsValid() {
			return identity.ErrInvalidRole
		}
		parsed = append(parsed, role)
	}
	filter.Roles = parsed
	return nil
}

func parseSearch(q url.Values, filter *identity.EmployeeFilter) {
	if search := q.Get("search"); search != "" {
		filter.Search = &search
	}
}

func parseActive(q url.Values, filter *identity.EmployeeFilter) error {
	activeStr := q.Get("active")
	if activeStr == "" {
		return nil
	}
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		return err
	}
	filter.Active = &active
	return nil
}

func parsePagination(q url.Values, filter *identity.EmployeeFilter) error {
	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return err
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
			return err
		}
		filter.Offset = offset
	}

	return nil
}

func parseSorting(q url.Values, filter *identity.EmployeeFilter) error {
	if orderBy := q.Get("order_by"); orderBy != "" {
		parsed, err := identity.ParseEmployeeOrderBy(orderBy)
		if err != nil {
			return err
		}
		filter.OrderBy = parsed
	}

	if sortOrder := q.Get("sort_order"); sortOrder != "" {
		parsed, err := shared.ParseSortOrder(sortOrder)
		if err != nil {
			return err
		}
		filter.SortOrder = parsed
	}

	return nil
}

func toGetListResponse(result *employeeUC.ListOutput) getListResponse {
	items := make([]employeeListItemResponse, 0, len(result.Items))

	for _, item := range result.Items {
		var locID *string
		if item.LocationID != nil {
			s := item.LocationID.String()
			locID = &s
		}

		items = append(items, employeeListItemResponse{
			ID:         item.ID.String(),
			Phone:      item.Phone,
			FirstName:  item.FirstName,
			LastName:   item.LastName,
			AvatarURL:  item.AvatarURL,
			LocationID: locID,
			Role:       string(item.Role),
			Permissions: permissionsResponse{
				CanProvideServices:        item.Permissions.CanProvideServices,
				CanManageLocationSchedule: item.Permissions.CanManageLocationSchedule,
			},
			Active:    item.Active,
			CreatedAt: item.CreatedAt,
		})
	}

	return getListResponse{
		Employees: items,
		Total:     result.Total,
	}
}
