// nolint: gosec
package appenv

const (
	DevSMSBotToken = "dev_sms_bot_token"
	DevSMSChatID   = "dev_sms_chat_id"

	SMSClientLogin          = "sms_client_login"
	SMSClientPassword       = "sms_client_password"
	RegisterConfirmationURL = "register_confirmation_url"

	SwaggerHost   = "swagger_host"
	SwaggerScheme = "swagger_scheme"

	WhatsAppAPIURL     = "whatsapp_api_url"
	WhatsAppMediaURL   = "whatsapp_media_url"
	WhatsAppIDInstance = "whatsapp_api_id_instance"
	WhatsAppAPIToken   = "whatsapp_api_token"

	JWTSecret       = "jwt_secret"
	JWTIssuer       = "jwt_issuer"
	AccessTokenTTL  = "access_token_ttl"
	RefreshTokenTTL = "refresh_token_ttl"
	RegistrationTTL = "registration_ttl"

	BagsyConfirmTTL = "bagsy_confirm_ttl"

	AwsRegion          = "aws_region"
	AwsS3BucketName    = "aws_s3_bucket_name"
	AwsS3Endpoint      = "aws_s3_endpoint"
	AwsSecretAccessKey = "aws_secret_access_key"
	AwsAccessKeyID     = "aws_access_key_id"

	MediaTTL             = "media_ttl"
	PointMediaMaxCount   = "point_media_max_count"
	PointCategoriesTTL   = "point_categories_ttl"
)
