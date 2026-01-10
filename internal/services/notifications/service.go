package notifications

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/log"
	"github.com/cockroachdb/errors"
)

type smsClient interface {
	Send(ctx context.Context, phone, message string) error
}

type whatsappClient interface {
	SendMessage(_ context.Context, phoneNumber, message string) error
}

type Service struct {
	smsClient              smsClient
	whatsApp               whatsappClient
	registrationConfirmURL string
}

func NewService(
	smsClient smsClient,
	whatsApp whatsappClient,
	registrationConfirmURL string,
) *Service {
	return &Service{
		smsClient:              smsClient,
		whatsApp:               whatsApp,
		registrationConfirmURL: registrationConfirmURL,
	}
}

func (s *Service) SendManagementAuthConfirmationCode(ctx context.Context, phone, code string) error {
	// TODO: format message (markdown)
	message := fmt.Sprintf("%s - Ваш код подтверждения регистрации в системе Bagsy", code)
	return s.send(ctx, phone, message)
}

func (s *Service) SendStaffRegistrationLink(ctx context.Context, phone, token string) error {
	link := fmt.Sprintf("%s/%s", s.registrationConfirmURL, token)
	// TODO: format message (markdown)
	message := fmt.Sprintf("Добро пожаловать в Bagsy! Завершите регистрацию по ссылке: %s", link)
	return s.send(ctx, phone, message)
}

func (s *Service) SendBagsyConfirmCode(ctx context.Context, phone, code string) error {
	// TODO: format message, add link, name of service etc. (markdown)
	message := fmt.Sprintf("%s - Ваш код подтверждения на запись", code)
	return s.send(ctx, phone, message)
}

func (s *Service) SendPasswordChangeLink(ctx context.Context, phone, token string) error {
	link := fmt.Sprintf("%s/%s", s.registrationConfirmURL, token)
	// TODO: format message (markdown)
	message := fmt.Sprintf("Для смены пароля в Bagsy следуйте по данной ссылке: %s", link)
	return s.send(ctx, phone, message)
}

func (s *Service) send(ctx context.Context, phone, message string) error {
	// Пытаемся отправить через WhatsApp
	err := s.whatsApp.SendMessage(ctx, phone, message)
	if err != nil {
		log.Warnf(ctx, "Failed to send message by whatsapp: %v", err)

		// Fallback на SMS
		err = s.smsClient.Send(ctx, phone, message)
		if err != nil {
			// Преобразуем SMS ошибку → доменную
			return s.mapSMSError(err)
		}
	}
	return nil
}

// mapSMSError преобразует ошибки SMS клиента → доменные ошибки
func (s *Service) mapSMSError(err error) error {
	// Validation errors
	if errors.Is(err, sms.ErrEmptyPhone) || errors.Is(err, sms.ErrEmptyMessage) {
		return domainErr.NewInvalidInputError("invalid notification data", err)
	}
	if errors.Is(err, sms.ErrInvalidPhone) {
		return domainErr.NewInvalidInputError("invalid phone number format", err)
	}

	// API errors - конвертируем в соответствующие доменные
	if errors.Is(err, sms.ErrAuthFailed) {
		return domainErr.NewUnauthorizedError("SMS service authentication failed")
	}
	if errors.Is(err, sms.ErrNoFunds) {
		return domainErr.NewInternalError("insufficient funds to send notification", err)
	}
	if errors.Is(err, sms.ErrIPBlocked) || errors.Is(err, sms.ErrForbidden) {
		return domainErr.NewInternalError("SMS service blocked the request", err)
	}
	if errors.Is(err, sms.ErrUndeliverable) {
		return domainErr.NewInternalError("notification cannot be delivered", err)
	}
	if errors.Is(err, sms.ErrTooManyRequests) || errors.Is(err, sms.ErrSpam) {
		return domainErr.NewInternalError("too many notification requests, try again later", err)
	}

	// Общая ошибка для всех остальных случаев
	return domainErr.NewInternalError("failed to send notification", err)
}

// mapWhatsAppError преобразует ошибки WhatsApp клиента → доменные ошибки
// (на случай если вам понадобится обрабатывать WhatsApp ошибки отдельно)
// nolint:unused // Функция создана для будущего использования
func (s *Service) mapWhatsAppError(err error) error {
	// Validation errors
	if errors.Is(err, whatsapp.ErrEmptyPhone) || errors.Is(err, whatsapp.ErrEmptyMessage) {
		return domainErr.NewInvalidInputError("invalid notification data", err)
	}
	if errors.Is(err, whatsapp.ErrEmptyFile) {
		return domainErr.NewInvalidInputError("file is required", err)
	}

	// API errors
	if errors.Is(err, whatsapp.ErrUnauthorized) {
		return domainErr.NewUnauthorizedError("WhatsApp API authentication failed")
	}
	if errors.Is(err, whatsapp.ErrInstanceOffline) {
		return domainErr.NewInternalError("WhatsApp instance is offline", err)
	}
	if errors.Is(err, whatsapp.ErrRateLimited) {
		return domainErr.NewInternalError("WhatsApp rate limit reached, try again later", err)
	}
	if errors.Is(err, whatsapp.ErrEmptyResponse) || errors.Is(err, whatsapp.ErrSendFailed) {
		return domainErr.NewInternalError("failed to send WhatsApp notification", err)
	}

	// Общая ошибка
	return domainErr.NewInternalError("failed to send notification", err)
}
