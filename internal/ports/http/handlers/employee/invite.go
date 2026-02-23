package employee

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	coreHTTP "github.com/Rasikrr/core/http"
)

// sendInvite handles POST /api/v1/employees/invite.
//
// @Summary      Пригласить сотрудника (шаг 1/2)
// @Description  Создаёт приглашение и отправляет ссылку для регистрации на номер телефона сотрудника через WhatsApp (с fallback на SMS). Требует авторизации. Owner может приглашать Manager и Staff, Manager — только Staff.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      sendInviteRequest  true  "Данные приглашения"
// @Success      201   {object}  sendInviteResponse
// @Failure      400   {object}  httputil.errorResponse  "Невалидные данные (телефон, роль, имя)"
// @Failure      401   {object}  httputil.errorResponse  "Не авторизован"
// @Failure      403   {object}  httputil.errorResponse  "Нет прав или подписка приостановлена"
// @Failure      409   {object}  httputil.errorResponse  "Сотрудник с таким телефоном уже существует"
// @Failure      429   {object}  httputil.errorResponse  "Приглашение уже отправлено (cooldown 60с)"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/employees/invite [post]
func (h *Handler) sendInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req sendInviteRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.inviteUseCase.SendInvite(ctx, orgCtx, uc.SendInviteInput{
		Phone:      req.Phone,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		LocationID: req.LocationID,
		Role:       req.Role,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, sendInviteResponse{
		Message:   "invite_sent",
		Phone:     out.Phone,
		ExpiresIn: out.ExpiresIn,
	}, http.StatusCreated)
}

// confirmInvite handles POST /api/v1/employees/invite/confirm.
//
// @Summary      Подтвердить приглашение (шаг 2/2)
// @Description  Сотрудник переходит по ссылке из приглашения, задаёт пароль. Создаётся аккаунт сотрудника, возвращаются JWT-токены. Не требует авторизации.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        body  body      confirmInviteRequest  true  "Токен приглашения и пароль"
// @Success      200   {object}  confirmInviteResponse
// @Failure      400   {object}  httputil.errorResponse  "Невалидные данные"
// @Failure      404   {object}  httputil.errorResponse  "Токен не найден"
// @Failure      409   {object}  httputil.errorResponse  "Телефон уже зарегистрирован"
// @Failure      410   {object}  httputil.errorResponse  "Токен истёк"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/employees/invite/confirm [post]
func (h *Handler) confirmInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req confirmInviteRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.inviteUseCase.ConfirmInvite(ctx, uc.ConfirmInviteInput{
		Token:    req.Token,
		Password: req.Password,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, confirmInviteResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}

// resendInvite handles POST /api/v1/employees/invite/resend.
//
// @Summary      Переотправить приглашение
// @Description  Генерирует новый токен и повторно отправляет ссылку для регистрации. Требует авторизации. Действует cooldown 60 секунд между отправками.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      resendInviteRequest  true  "Номер телефона приглашённого"
// @Success      200   {object}  resendInviteResponse
// @Failure      400   {object}  httputil.errorResponse  "Невалидный телефон"
// @Failure      401   {object}  httputil.errorResponse  "Не авторизован"
// @Failure      403   {object}  httputil.errorResponse  "Нет прав (другая организация)"
// @Failure      404   {object}  httputil.errorResponse  "Приглашение не найдено"
// @Failure      429   {object}  httputil.errorResponse  "Cooldown ещё не прошёл"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/employees/invite/resend [post]
func (h *Handler) resendInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req resendInviteRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.inviteUseCase.ResendInvite(ctx, orgCtx, uc.ResendInviteInput{
		Phone: req.Phone,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, employeeErrors)
		return
	}

	coreHTTP.SendData(ctx, w, resendInviteResponse{
		Message:    "invite_resent",
		Phone:      out.Phone,
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)
}
