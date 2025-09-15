// nolint: unused
package sms

type responseFormat uint8

const (
	responseFormatString responseFormat = iota
	responseFormatNumber
	responseFormatXML
	responseFormatJSON
)

type smsStatus int8

const (
	smsStatusNotFound smsStatus = iota - 3
	smsStatusStopped
	smsStatusPending
	smsStatusPassedToOperator
	smsSatusDelivered
	smsStatusChecked
	smsStatusExpired
	smsStatusClicked
	smsStatusImpossibleToDeliver = iota + 12
	smsStatusInvalidNumber       = iota + 13
	smsStatusForbidden
	smsStatusInsufficientFunds
	smsStatusUnavailableNumber
)

var (
	errSmsStatuses = []smsStatus{
		smsStatusNotFound,
		smsStatusStopped,
		smsStatusExpired,
		smsStatusImpossibleToDeliver,
		smsStatusInvalidNumber,
		smsStatusForbidden,
		smsStatusInsufficientFunds,
		smsStatusUnavailableNumber,
	}
)

func (s smsStatus) OneOf(statuses ...smsStatus) bool {
	for _, status := range statuses {
		if s == status {
			return true
		}
	}
	return false
}

type errCodes uint8

const (
	errParams           errCodes = iota + 1 // 1. Ошибка в параметрах
	errAuth                                 // 2. Неверный логин или пароль / IP-адрес не в списке
	errNoFunds                              // 3. Недостаточно средств на счете
	errIPBlocked                            // 4. IP-адрес временно заблокирован
	errDateFormat                           // 5. Неверный формат даты
	errMessageForbidden                     // 6. Сообщение запрещено (по тексту или имени отправителя)
	errPhoneFormat                          // 7. Неверный формат номера телефона
	errUndeliverable                        // 8. Сообщение не может быть доставлено
	errTooManyRequests                      // 9. Отправка более одного одинакового запроса / слишком много concurrent requests
)
