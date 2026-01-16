// nolint: godot
package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getUsers godoc
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с offset-based пагинацией и учетом прав доступа.
// @Description - NetManager/SelfOwner: пользователи только своей сети
// @Description - Manager: пользователи только своей точки
// @Description - Staff/User: доступ запрещен (403)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param point_code query string false "Код точки для фильтрации"
// @Param network_code query string false "Код сети для фильтрации"
// @Param role query []string false "Фильтр по ролям" collectionFormat(multi) Enums(staff,manager,net_manager,self_owner)
// @Param phone_search query string false "Поиск по части или полному номеру телефона (поиск в начале, середине и конце)"
// @Param limit query int false "Количество записей на странице (default: 20, max: 100)"
// @Param offset query int false "Смещение для пагинации (default: 0)"
// @Param order_by query string false "Поле для сортировки" Enums(phone,name,surname,point_code,network_code,created_at,updated_at) default(created_at)
// @Param sort_order query string false "Направление сортировки" Enums(asc,desc) default(asc)
// @Success 200 {object} getUsersResponse "Список пользователей с пагинацией"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для доступа к ресурсу"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/staff [get]
func (c *Controller) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getUsersRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	filter, err := req.toFilter()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	result, err := c.userService.GetListByFilter(ctx, filter)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toGetUsersResponse(result), http.StatusOK)
}
