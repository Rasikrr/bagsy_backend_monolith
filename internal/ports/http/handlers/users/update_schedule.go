// nolint: godot
package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// updateSchedule godoc
// @Summary Обновление расписания пользователя
// @Description Позволяет обновить расписание работы (7 дней). Время передается в формате Almaty timezone и конвертируется в UTC
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body updateScheduleRequest true "Расписание на 7 дней"
// @Success 200 {object} response.EmptySuccessResponse "Расписание обновлено"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 404 {object} errors.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/me/schedule [put]
func (c *Controller) updateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actor, err := actor.GetActor(ctx)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	var req updateScheduleRequest
	if err = request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	schedules, err := req.toDomain()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err = c.userService.UpdateSchedule(ctx, actor.Phone(), schedules)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("schedule updated"), http.StatusOK)
}
