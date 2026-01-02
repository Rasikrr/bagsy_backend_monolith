package sms

import (
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

// SendRequest структура запроса для отправки SMS
type sendRequest struct {
	Login    string         `json:"login"`              // Логин клиента
	Password string         `json:"psw"`                // Пароль или API ключ
	Phones   string         `json:"phones"`             // Номер телефона получателя
	Message  string         `json:"mes"`                // Текст сообщения
	Format   ResponseFormat `json:"fmt"`                // Формат ответа
	Charset  string         `json:"charset"`            // Кодировка (utf-8)
	Translit int            `json:"translit,omitempty"` // Транслит (0 или 1)
	Sender   string         `json:"sender,omitempty"`   // Имя отправителя
}

// sendResponse структура ответа на отправку SMS
type sendResponse struct {
	ID        int       `json:"id,omitempty"`         // ID сообщения
	Count     int       `json:"cnt,omitempty"`        // Количество SMS
	Cost      string    `json:"cost,omitempty"`       // Стоимость
	Balance   string    `json:"balance,omitempty"`    // Баланс после отправки
	Error     string    `json:"error,omitempty"`      // Описание ошибки
	ErrorCode ErrorCode `json:"error_code,omitempty"` // Код ошибки
}

// HasError проверяет наличие ошибки в ответе
func (r *sendResponse) HasError() bool {
	return r.Error != "" || r.ErrorCode != 0
}

// GetError возвращает ошибку из ответа и конвертирует в доменную
func (r *sendResponse) GetError() error {
	if !r.HasError() {
		return nil
	}

	// Конвертируем коды ошибок API в доменные ошибки
	switch r.ErrorCode {
	case ErrorCodeAuth:
		return domainErr.ErrSMSAuthFailed.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeNoFunds:
		return domainErr.ErrSMSNoFunds.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeIPBlocked:
		return domainErr.ErrSMSIPBlocked.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeMessageForbidden:
		return domainErr.ErrSMSForbidden.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodePhoneFormat:
		return domainErr.ErrSMSInvalidPhone.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeUndeliverable:
		return domainErr.ErrSMSUndeliverable.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeTooManyRequests:
		return domainErr.ErrSMSTooManyRequests.WithError(nil).WithDetail("api_error", r.Error)
	default:
		return domainErr.ErrSMSSendFailed.WithError(nil).WithDetail("api_error", r.Error).
			WithDetail("error_code", int(r.ErrorCode))
	}
}

// StatusResponse структура ответа на проверку статуса
type StatusResponse struct {
	Status        Status    `json:"status"`                   // Статус сообщения
	LastDate      string    `json:"last_date,omitempty"`      // Дата последнего изменения
	LastTimestamp int64     `json:"last_timestamp,omitempty"` // Unix timestamp
	Flag          int       `json:"flag,omitempty"`           // Дополнительный флаг
	Error         string    `json:"error,omitempty"`          // Ошибка
	ErrorCode     ErrorCode `json:"error_code,omitempty"`     // Код ошибки
}

// HasError проверяет наличие ошибки в ответе
func (r *StatusResponse) HasError() bool {
	return r.Error != "" || r.ErrorCode != 0
}

// GetError возвращает ошибку из ответа и конвертирует в доменную
func (r *StatusResponse) GetError() error {
	if !r.HasError() {
		return nil
	}

	// Конвертируем коды ошибок API в доменные ошибки
	switch r.ErrorCode {
	case ErrorCodeAuth:
		return domainErr.ErrSMSAuthFailed.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeNoFunds:
		return domainErr.ErrSMSNoFunds.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodeIPBlocked:
		return domainErr.ErrSMSIPBlocked.WithError(nil).WithDetail("api_error", r.Error)
	case ErrorCodePhoneFormat:
		return domainErr.ErrSMSInvalidPhone.WithError(nil).WithDetail("api_error", r.Error)
	default:
		return domainErr.ErrSMSSendFailed.WithError(nil).WithDetail("api_error", r.Error).WithDetail("error_code", int(r.ErrorCode))
	}
}
