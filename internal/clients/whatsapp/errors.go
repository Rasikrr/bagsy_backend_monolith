package whatsapp

import "errors"

var (
	ErrWhatsAppPhoneRequired    = errors.New("phone number is required")
	ErrWhatsAppMessageRequired  = errors.New("message is required")
	ErrWhatsAppFileRequired     = errors.New("file is required")
	ErrWhatsAppEmptyResponse    = errors.New("empty response from API")
	ErrWhatsAppInstanceOffline  = errors.New("instance is offline")
	ErrWhatsAppInvalidPhone     = errors.New("invalid phone number format")
	ErrWhatsAppSendFailed       = errors.New("failed to send message")
	ErrWhatsAppAccountNotFound  = errors.New("whatsapp account not found")
	ErrWhatsAppUnauthorized     = errors.New("unauthorized access to whatsapp API")
	ErrWhatsAppRateLimitReached = errors.New("rate limit reached")
)
