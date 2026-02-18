package messenger

import "fmt"

const (
	otpTemplate           = "Ваш код подтверждения Bagsy: %s"
	passwordResetTemplate = "Для сброса пароля в Bagsy перейдите по ссылке: %s"
)

func formatOTPMessage(code string) string {
	return fmt.Sprintf(otpTemplate, code)
}

func formatPasswordResetMessage(link string) string {
	return fmt.Sprintf(passwordResetTemplate, link)
}
