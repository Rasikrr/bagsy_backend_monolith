package errors

// SMS errors
var (
	// Validation errors
	ErrSMSEmptyMessage = NewInvalidInputError("message cannot be empty", nil)
	ErrSMSEmptyPhone   = NewInvalidInputError("phone number cannot be empty", nil)
	ErrSMSInvalidPhone = NewInvalidInputError("invalid phone number format", nil)

	// API errors
	ErrSMSAuthFailed      = NewUnauthorizedError("SMS service authentication failed")
	ErrSMSNoFunds         = NewInternalError("insufficient funds on SMS service account", nil)
	ErrSMSIPBlocked       = NewInternalError("IP address blocked by SMS service", nil)
	ErrSMSForbidden       = NewInternalError("message forbidden by SMS service", nil)
	ErrSMSUndeliverable   = NewInternalError("SMS cannot be delivered", nil)
	ErrSMSTooManyRequests = NewInternalError("too many SMS requests", nil)
	ErrSMSSpam            = NewConflictError("spam detected, try again later", nil)
	ErrSMSSendFailed      = NewInternalError("failed to send SMS", nil)
)

// WhatsApp errors
var (
	// Validation errors
	ErrWhatsAppEmptyPhone   = NewInvalidInputError("phone number is required", nil)
	ErrWhatsAppEmptyMessage = NewInvalidInputError("message is required", nil)
	ErrWhatsAppEmptyFile    = NewInvalidInputError("file is required", nil)

	// API errors
	ErrWhatsAppSendFailed      = NewInternalError("failed to send WhatsApp message", nil)
	ErrWhatsAppEmptyResponse   = NewInternalError("empty response from WhatsApp API", nil)
	ErrWhatsAppInstanceOffline = NewInternalError("WhatsApp instance is offline", nil)
	ErrWhatsAppUnauthorized    = NewUnauthorizedError("unauthorized access to WhatsApp API")
	ErrWhatsAppRateLimited     = NewInternalError("WhatsApp rate limit reached", nil)
)
