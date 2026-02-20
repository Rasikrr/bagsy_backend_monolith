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
	ErrEmployeeInactive = errors.New("employee account is inactive")
)

var (
	ErrInviteTokenNotFound = errors.New("invite token not found")
	ErrInviteTokenExpired  = errors.New("invite token has expired")
	ErrInviteAlreadyExists = errors.New("invite already exists for this phone")
)

var (
	ErrUnknownTokenPurpose = errors.New("unknown token purpose")
	ErrActionTokenNotFound = errors.New("action token not found")
)
