package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getCustomers godoc
// @Summary Получить список клиентов
// @Description Возвращает список клиентов (users с role='user'), обслуживавшихся в точках, с offset-based пагинацией и учетом прав доступа.
// @Description - Staff: только клиенты которых он обслуживал (можно фильтровать по своей точке)
// @Description - Manager: все клиенты своей точки (можно указать point_code=своя_точка, можно фильтровать по телефону мастера)
// @Description - NetManager/SelfOwner: клиенты сети (можно фильтровать по конкретным точкам или получить всех, можно фильтровать по телефону мастера)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param phone_search query string false "Поиск по части или полному номеру телефона клиента"
// @Param staff_phone query string false "Фильтр по телефону мастера (для Manager и выше)"
// @Param point_code query []string false "Фильтр по кодам точек (можно указать несколько)" collectionFormat(multi)
// @Param limit query int false "Количество записей на странице (default: 20, max: 100)"
// @Param offset query int false "Смещение для пагинации (default: 0)"
// @Param order_by query string false "Поле для сортировки" Enums(phone,name,surname,created_at,updated_at) default(created_at)
// @Param sort_order query string false "Направление сортировки" Enums(asc,desc) default(asc)
// @Success 200 {object} getCustomersResponse "Список клиентов с пагинацией"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для доступа к ресурсу"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/customers [get]
func (c *Controller) getCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getCustomersRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	filter, err := req.toFilter()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	result, err := c.userService.GetCustomers(ctx, filter)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toGetCustomersResponse(result), http.StatusOK)
}
