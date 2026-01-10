package errors

// NOTE: SMS and WhatsApp errors have been moved to their respective client packages:
// - SMS errors: internal/clients/sms/errors.go
// - WhatsApp errors: internal/clients/whatsapp/errors.go
//
// This follows Clean Architecture principles - domain should not know about infrastructure details.
// Services should convert infrastructure errors to domain errors when needed.
