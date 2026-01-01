package sms

// ResponseFormat определяет формат ответа от API
type ResponseFormat int

const (
	ResponseFormatPlainText ResponseFormat = iota // 0 - обычная строка
	ResponseFormatCSV                             // 1 - строка с разделителями-запятыми
	ResponseFormatXML                             // 2 - XML-документ
	ResponseFormatJSON                            // 3 - JSON
)

// ErrorCode коды ошибок от SMSC.KZ API
type ErrorCode int

const (
	ErrorCodeParams           ErrorCode = 1 // Ошибка в параметрах
	ErrorCodeAuth             ErrorCode = 2 // Неверный логин или пароль / IP не в списке
	ErrorCodeNoFunds          ErrorCode = 3 // Недостаточно средств
	ErrorCodeIPBlocked        ErrorCode = 4 // IP временно заблокирован
	ErrorCodeDateFormat       ErrorCode = 5 // Неверный формат даты
	ErrorCodeMessageForbidden ErrorCode = 6 // Сообщение запрещено
	ErrorCodePhoneFormat      ErrorCode = 7 // Неверный формат телефона
	ErrorCodeUndeliverable    ErrorCode = 8 // Невозможно доставить
	ErrorCodeTooManyRequests  ErrorCode = 9 // Слишком много запросов
)

// SMSStatus статусы SMS сообщения
type SMSStatus int

const (
	SMSStatusNotFound          SMSStatus = -3 // Не найдено
	SMSStatusStopped           SMSStatus = -2 // Остановлено
	SMSStatusPending           SMSStatus = -1 // Ожидает отправки
	SMSStatusPassedToOperator  SMSStatus = 0  // Передано оператору
	SMSStatusDelivered         SMSStatus = 1  // Доставлено
	SMSStatusRead              SMSStatus = 2  // Прочитано
	SMSStatusExpired           SMSStatus = 3  // Срок истек
	SMSStatusClicked           SMSStatus = 4  // Нажат переход по ссылке
	SMSStatusImpossibleDeliver SMSStatus = 20 // Невозможно доставить
	SMSStatusInvalidNumber     SMSStatus = 21 // Неправильный номер
	SMSStatusForbidden         SMSStatus = 22 // Запрещено
	SMSStatusInsufficientFunds SMSStatus = 23 // Недостаточно средств
	SMSStatusUnavailableNumber SMSStatus = 24 // Номер недоступен
)

// IsError проверяет, является ли статус ошибочным
func (s SMSStatus) IsError() bool {
	errorStatuses := []SMSStatus{
		SMSStatusNotFound,
		SMSStatusStopped,
		SMSStatusExpired,
		SMSStatusImpossibleDeliver,
		SMSStatusInvalidNumber,
		SMSStatusForbidden,
		SMSStatusInsufficientFunds,
		SMSStatusUnavailableNumber,
	}

	for _, status := range errorStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsSuccess проверяет, успешно ли доставлено сообщение
func (s SMSStatus) IsSuccess() bool {
	return s == SMSStatusDelivered || s == SMSStatusRead || s == SMSStatusClicked
}
