// nolint: gosec
package appenv

const (
	DevSMSBotToken = "dev_sms_bot_token"
	DevSMSChatID   = "dev_sms_chat_id"

	AuthCodeTTL             = "auth_code_ttl"
	BagsiesCodeTTL          = "bagsies_code_ttl"
	SMSSpamTTL              = "sms_spam_ttl"
	SMSClientLogin          = "sms_client_login"
	SMSClientPassword       = "sms_client_password"
	RegisterConfirmationURL = "register_confirmation_url"

	SwaggerHost   = "swagger_host"
	SwaggerScheme = "swagger_scheme"

	WhatsAppAPIURL     = "whatsapp_api_url"
	WhatsAppMediaURL   = "whatsapp_media_url"
	WhatsAppIDInstance = "whatsapp_api_id_instance"
	WhatsAppAPIToken   = "whatsapp_api_token"

	InactiveUserTTL         = "inactive_user_ttl"
	InactiveUserJobSchedule = "inactive_user_job_schedule"

	JWTSecret            = "jwt_secret"
	JWTIssuer            = "jwt_issuer"
	AccessTokenTTL       = "access_token_ttl"
	RefreshTokenTTL      = "refresh_token_ttl"
	RegistrationTokenTTL = "registration_token_ttl"
)
