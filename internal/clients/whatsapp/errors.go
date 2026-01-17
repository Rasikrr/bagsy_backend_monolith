package whatsapp

import "github.com/cockroachdb/errors"

// Validation errors
var (
	ErrEmptyPhone   = errors.New("whatsapp: phone number is required")
	ErrEmptyMessage = errors.New("whatsapp: message is required")
	ErrEmptyFile    = errors.New("whatsapp: file is required")
	ErrEmptyChatID  = errors.New("whatsapp: chat ID is required")
	ErrEmptyMsgID   = errors.New("whatsapp: message ID is required")
	ErrEmptyContact = errors.New("whatsapp: contact phone number is required")
)

// API errors
var (
	ErrSendFailed      = errors.New("whatsapp: failed to send message")
	ErrEmptyResponse   = errors.New("whatsapp: empty response from API")
	ErrInstanceOffline = errors.New("whatsapp: instance is offline")
	ErrUnauthorized    = errors.New("whatsapp: unauthorized access to API")
	ErrRateLimited     = errors.New("whatsapp: rate limit reached")
)

// Internal errors
var (
	ErrGetStateFailed    = errors.New("whatsapp: failed to get state instance")
	ErrGetSettingsFailed = errors.New("whatsapp: failed to get settings")
	ErrSetSettingsFailed = errors.New("whatsapp: failed to set settings")
	ErrRebootFailed      = errors.New("whatsapp: failed to reboot instance")
	ErrLogoutFailed      = errors.New("whatsapp: failed to logout")
	ErrDownloadFailed    = errors.New("whatsapp: failed to download file")
	ErrUnmarshalFailed   = errors.New("whatsapp: failed to unmarshal response")
)
