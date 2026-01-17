package sms

import "github.com/cockroachdb/errors"

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

// GetError возвращает ошибку из ответа
func (r *sendResponse) GetError() error {
	if !r.HasError() {
		return nil
	}

	// Конвертируем коды ошибок API в SMS ошибки
	var baseErr error
	switch r.ErrorCode {
	case ErrorCodeAuth:
		baseErr = ErrAuthFailed
	case ErrorCodeNoFunds:
		baseErr = ErrNoFunds
	case ErrorCodeIPBlocked:
		baseErr = ErrIPBlocked
	case ErrorCodeMessageForbidden:
		baseErr = ErrForbidden
	case ErrorCodePhoneFormat:
		baseErr = ErrInvalidPhone
	case ErrorCodeUndeliverable:
		baseErr = ErrUndeliverable
	case ErrorCodeTooManyRequests:
		baseErr = ErrTooManyRequests
	default:
		baseErr = ErrSendFailed
	}

	return errors.Wrapf(baseErr, "api_error=%s, error_code=%d", r.Error, int(r.ErrorCode))
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

// GetError возвращает ошибку из ответа
func (r *StatusResponse) GetError() error {
	if !r.HasError() {
		return nil
	}

	// Конвертируем коды ошибок API в SMS ошибки
	var baseErr error
	switch r.ErrorCode {
	case ErrorCodeAuth:
		baseErr = ErrAuthFailed
	case ErrorCodeNoFunds:
		baseErr = ErrNoFunds
	case ErrorCodeIPBlocked:
		baseErr = ErrIPBlocked
	case ErrorCodePhoneFormat:
		baseErr = ErrInvalidPhone
	default:
		baseErr = ErrSendFailed
	}

	return errors.Wrapf(baseErr, "api_error=%s, error_code=%d", r.Error, int(r.ErrorCode))
}
