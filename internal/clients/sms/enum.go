package sms

// ResponseFormat определяет формат ответа от API
type ResponseFormat int

const (
	ResponseFormatPlainText ResponseFormat = iota // 0 - обычная строка
	ResponseFormatCSV                             // 1 - строка с разделителями-запятыми
	ResponseFormatXML                             // 2 - XML-документ
	ResponseFormatJSON                            // 3 - JSON
)

func (r ResponseFormat) Int() int {
	return int(r)
}

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

// Status статусы SMS сообщения
type Status int

const (
	StatusNotFound          Status = -3 // Не найдено
	StatusStopped           Status = -2 // Остановлено
	StatusPending           Status = -1 // Ожидает отправки
	StatusPassedToOperator  Status = 0  // Передано оператору
	StatusDelivered         Status = 1  // Доставлено
	StatusRead              Status = 2  // Прочитано
	StatusExpired           Status = 3  // Срок истек
	StatusClicked           Status = 4  // Нажат переход по ссылке
	StatusImpossibleDeliver Status = 20 // Невозможно доставить
	StatusInvalidNumber     Status = 21 // Неправильный номер
	StatusForbidden         Status = 22 // Запрещено
	StatusInsufficientFunds Status = 23 // Недостаточно средств
	StatusUnavailableNumber Status = 24 // Номер недоступен
)

// IsError проверяет, является ли статус ошибочным
func (s Status) IsError() bool {
	errorStatuses := []Status{
		StatusNotFound,
		StatusStopped,
		StatusExpired,
		StatusImpossibleDeliver,
		StatusInvalidNumber,
		StatusForbidden,
		StatusInsufficientFunds,
		StatusUnavailableNumber,
	}

	for _, status := range errorStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsSuccess проверяет, успешно ли доставлено сообщение
func (s Status) IsSuccess() bool {
	return s == StatusDelivered || s == StatusRead || s == StatusClicked
}
