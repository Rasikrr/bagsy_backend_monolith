package auth

import "errors"

var (
	ErrOTPExpired     = errors.New("OTP code has expired")
	ErrOTPInvalid     = errors.New("OTP code is invalid")
	ErrOTPAlreadySent = errors.New("OTP code already sent recently")
)
