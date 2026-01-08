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
	ErrSMSRequestFailed   = NewInternalError("failed to send request to SMS client", nil)
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

// S3 errors
var (
	// Validation errors
	ErrS3EmptyRegion    = NewInvalidInputError("AWS region is required", nil)
	ErrS3EmptyAccessKey = NewInvalidInputError("AWS access key is required", nil)
	ErrS3EmptySecretKey = NewInvalidInputError("AWS secret key is required", nil)
	ErrS3EmptyBucket    = NewInvalidInputError("S3 bucket name is required", nil)
	ErrS3EmptyKey       = NewInvalidInputError("S3 object key is required", nil)
	ErrS3EmptyData      = NewInvalidInputError("data to upload cannot be empty", nil)

	// API errors
	ErrS3ConfigFailed   = NewInternalError("failed to load AWS configuration", nil)
	ErrS3UploadFailed   = NewInternalError("failed to upload file to S3", nil)
	ErrS3DownloadFailed = NewInternalError("failed to download file from S3", nil)
	ErrS3DeleteFailed   = NewInternalError("failed to delete file from S3", nil)
	ErrS3ListFailed     = NewInternalError("failed to list S3 objects", nil)
	ErrS3EmptyLocation  = NewInternalError("empty location returned from S3", nil)
)
