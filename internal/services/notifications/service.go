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

type messenger interface {
	SendMessage(_ context.Context, phoneNumber, message string) error
}

type Service struct {
	smsClient              messenger
	whatsApp               messenger
	registrationConfirmURL string
}

func NewService(
	smsClient messenger,
	whatsApp messenger,
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
	whatsappErr := s.whatsApp.SendMessage(ctx, phone, message)
	if whatsappErr != nil {
		log.Warnf(ctx, "Failed to send message by whatsapp: %v", whatsappErr)

		// Fallback на SMS
		smsErr := s.smsClient.SendMessage(ctx, phone, message)
		if smsErr != nil {
			// Обе попытки отправки провалились - возвращаем доменную ошибку с деталями обеих ошибок
			mappedErr := s.mapNotificationError(smsErr)
			// Добавляем информацию об ошибке WhatsApp в Details
			return mappedErr.
				WithDetail("whatsapp_error", whatsappErr.Error()).
				WithDetail("sms_error", smsErr.Error()).
				WithDetail("fallback_attempted", true)
		}
		// SMS отправился успешно после fallback
		log.Infof(ctx, "Message sent via SMS fallback after WhatsApp failure")
	}
	return nil
}

// mapNotificationError преобразует ошибки notification клиентов (SMS/WhatsApp) → доменные ошибки
// nolint:gocognit,gocyclo,cyclop,funlen // Comprehensive error mapping requires explicit handling of all error cases
func (s *Service) mapNotificationError(err error) *domainErr.DomainError {
	if err == nil {
		return nil
	}

	// ========== VALIDATION ERRORS (400) ==========

	// Phone validation (SMS + WhatsApp)
	if errors.Is(err, sms.ErrEmptyPhone) || errors.Is(err, whatsapp.ErrEmptyPhone) {
		return domainErr.NewInvalidInputError("phone number is required", err)
	}
	if errors.Is(err, sms.ErrInvalidPhone) {
		return domainErr.NewInvalidInputError("invalid phone number format", err)
	}

	// Message validation (SMS + WhatsApp)
	if errors.Is(err, sms.ErrEmptyMessage) || errors.Is(err, whatsapp.ErrEmptyMessage) {
		return domainErr.NewInvalidInputError("message is required", err)
	}

	// WhatsApp specific validation
	if errors.Is(err, whatsapp.ErrEmptyFile) {
		return domainErr.NewInvalidInputError("file is required", err)
	}
	if errors.Is(err, whatsapp.ErrEmptyChatID) {
		return domainErr.NewInvalidInputError("chat ID is required", err)
	}
	if errors.Is(err, whatsapp.ErrEmptyContact) {
		return domainErr.NewInvalidInputError("contact is required", err)
	}

	// Message ID validation (SMS + WhatsApp)
	if errors.Is(err, sms.ErrInvalidMsgID) || errors.Is(err, whatsapp.ErrEmptyMsgID) {
		return domainErr.NewInvalidInputError("message ID is invalid", err)
	}

	// ========== AUTH ERRORS (401) ==========

	if errors.Is(err, sms.ErrAuthFailed) || errors.Is(err, whatsapp.ErrUnauthorized) {
		return domainErr.NewUnauthorizedError("notification service authentication failed")
	}

	// ========== BUSINESS/API ERRORS (500) ==========

	// Insufficient funds (SMS only)
	if errors.Is(err, sms.ErrNoFunds) {
		return domainErr.NewInternalError("insufficient funds to send notification", err)
	}

	// IP/Access restrictions (SMS only)
	if errors.Is(err, sms.ErrIPBlocked) {
		return domainErr.NewInternalError("IP address blocked by notification service", err)
	}
	if errors.Is(err, sms.ErrForbidden) {
		return domainErr.NewInternalError("notification service rejected the message", err)
	}

	// Delivery issues (SMS only)
	if errors.Is(err, sms.ErrUndeliverable) {
		return domainErr.NewInternalError("notification cannot be delivered", err)
	}

	// Rate limiting (SMS + WhatsApp)
	if errors.Is(err, sms.ErrTooManyRequests) || errors.Is(err, whatsapp.ErrRateLimited) {
		return domainErr.NewInternalError("too many notification requests, try again later", err)
	}

	// Spam detection (SMS only)
	if errors.Is(err, sms.ErrSpam) {
		return domainErr.NewInternalError("notification rejected as spam, try again later", err)
	}

	// WhatsApp instance offline
	if errors.Is(err, whatsapp.ErrInstanceOffline) {
		return domainErr.NewInternalError("WhatsApp instance is temporarily unavailable", err)
	}

	// Generic send failures (SMS + WhatsApp)
	if errors.Is(err, sms.ErrSendFailed) || errors.Is(err, whatsapp.ErrSendFailed) {
		return domainErr.NewInternalError("failed to send notification", err)
	}

	// Empty response (WhatsApp only)
	if errors.Is(err, whatsapp.ErrEmptyResponse) {
		return domainErr.NewInternalError("received empty response from notification service", err)
	}

	// ========== INTERNAL/NETWORK ERRORS (500) ==========

	// SMS internal errors
	if errors.Is(err, sms.ErrMarshalFailed) || errors.Is(err, sms.ErrUnmarshalFailed) {
		return domainErr.NewInternalError("notification service serialization error", err)
	}
	if errors.Is(err, sms.ErrCreateRequestFailed) || errors.Is(err, sms.ErrRequestFailed) {
		return domainErr.NewInternalError("failed to create notification request", err)
	}
	if errors.Is(err, sms.ErrHTTPRequestFailed) {
		return domainErr.NewInternalError("notification service network error", err)
	}
	if errors.Is(err, sms.ErrUnexpectedStatus) {
		return domainErr.NewInternalError("notification service returned unexpected status", err)
	}
	if errors.Is(err, sms.ErrReadBodyFailed) {
		return domainErr.NewInternalError("failed to read notification service response", err)
	}

	// WhatsApp internal errors
	if errors.Is(err, whatsapp.ErrUnmarshalFailed) {
		return domainErr.NewInternalError("WhatsApp service serialization error", err)
	}
	if errors.Is(err, whatsapp.ErrGetStateFailed) ||
		errors.Is(err, whatsapp.ErrGetSettingsFailed) ||
		errors.Is(err, whatsapp.ErrSetSettingsFailed) {
		return domainErr.NewInternalError("WhatsApp service configuration error", err)
	}
	if errors.Is(err, whatsapp.ErrRebootFailed) || errors.Is(err, whatsapp.ErrLogoutFailed) {
		return domainErr.NewInternalError("WhatsApp service operation error", err)
	}
	if errors.Is(err, whatsapp.ErrDownloadFailed) {
		return domainErr.NewInternalError("failed to download file from WhatsApp", err)
	}

	// ========== FALLBACK ==========

	// Для всех неизвестных ошибок
	return domainErr.NewInternalError("failed to send notification", err)
}
