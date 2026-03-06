// nolint: gosec
package messenger

import "fmt"

const (
	otpTemplate           = "Ваш код подтверждения Bagsy: %s"
	bookingOTPTemplate    = "Код подтверждения записи в Bagsy: %s"
	passwordResetTemplate = "Для сброса пароля в Bagsy перейдите по ссылке: %s"
	inviteTemplate        = "Вас пригласили в команду Bagsy! Для регистрации перейдите по ссылке: %s"

	customer24hReminderTemplate = "Напоминаем: завтра в %s у вас запись на %s (%s). Bagsy"
	customer1hReminderTemplate  = "Через час у вас запись на %s (%s). Bagsy"
	employee24hReminderTemplate = "Напоминание: завтра в %s запись клиента на %s (%s). Bagsy"
	employee1hReminderTemplate  = "Через час запись клиента на %s (%s). Bagsy"
)

func formatOTPMessage(code string) string {
	return fmt.Sprintf(otpTemplate, code)
}

func formatBookingOTPMessage(code string) string {
	return fmt.Sprintf(bookingOTPTemplate, code)
}

func formatPasswordResetMessage(link string) string {
	return fmt.Sprintf(passwordResetTemplate, link)
}

func formatInviteMessage(link string) string {
	return fmt.Sprintf(inviteTemplate, link)
}

func formatCustomer24hReminder(startAt, serviceName, locationName string) string {
	return fmt.Sprintf(customer24hReminderTemplate, startAt, serviceName, locationName)
}

func formatCustomer1hReminder(serviceName, locationName string) string {
	return fmt.Sprintf(customer1hReminderTemplate, serviceName, locationName)
}

func formatEmployee24hReminder(startAt, serviceName, locationName string) string {
	return fmt.Sprintf(employee24hReminderTemplate, startAt, serviceName, locationName)
}

func formatEmployee1hReminder(serviceName, locationName string) string {
	return fmt.Sprintf(employee1hReminderTemplate, serviceName, locationName)
}
