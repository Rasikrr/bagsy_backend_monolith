package auth

import "errors"

var (
	ErrOTPExpired      = errors.New("OTP code has expired")
	ErrOTPInvalid      = errors.New("OTP code is invalid")
	ErrOTPAlreadySent  = errors.New("OTP code already sent recently")
	ErrTooManyAttempts = errors.New("too many OTP verification attempts")
)

var (
	ErrPhoneAlreadyExists   = errors.New("phone already exists")
	ErrRegistrationExpired  = errors.New("registration has expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")
)

var (
	ErrResetTokenNotFound = errors.New("password reset token not found")
	ErrResetTokenExpired  = errors.New("password reset token has expired")
	ErrEmployeeInactive   = errors.New("employee account is inactive")
)
