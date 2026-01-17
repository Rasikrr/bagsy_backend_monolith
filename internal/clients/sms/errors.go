package sms

import "github.com/cockroachdb/errors"

// Validation errors
var (
	ErrEmptyMessage = errors.New("sms: message cannot be empty")
	ErrEmptyPhone   = errors.New("sms: phone number cannot be empty")
	ErrInvalidPhone = errors.New("sms: invalid phone number format")
	ErrInvalidMsgID = errors.New("sms: invalid message ID")
)

// API errors
var (
	ErrAuthFailed      = errors.New("sms: service authentication failed")
	ErrNoFunds         = errors.New("sms: insufficient funds on service account")
	ErrIPBlocked       = errors.New("sms: IP address blocked by service")
	ErrForbidden       = errors.New("sms: message forbidden by service")
	ErrUndeliverable   = errors.New("sms: SMS cannot be delivered")
	ErrTooManyRequests = errors.New("sms: too many requests")
	ErrSpam            = errors.New("sms: spam detected")
	ErrSendFailed      = errors.New("sms: failed to send SMS")
	ErrRequestFailed   = errors.New("sms: failed to send request to SMS client")
)

// Internal errors
var (
	ErrMarshalFailed       = errors.New("sms: failed to marshal request body")
	ErrCreateRequestFailed = errors.New("sms: failed to create request")
	ErrHTTPRequestFailed   = errors.New("sms: failed to execute HTTP request")
	ErrUnexpectedStatus    = errors.New("sms: unexpected API status code")
	ErrReadBodyFailed      = errors.New("sms: failed to read response body")
	ErrUnmarshalFailed     = errors.New("sms: failed to unmarshal response")
)
