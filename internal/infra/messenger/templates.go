package messenger

import "fmt"

const otpTemplate = "Ваш код подтверждения Bagsy: %s"

func formatOTPMessage(code string) string {
	return fmt.Sprintf(otpTemplate, code)
}
