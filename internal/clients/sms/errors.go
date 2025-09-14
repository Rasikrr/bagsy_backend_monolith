package sms

import "errors"

var (
	errEmptyMessage = errors.New("empty message")
	errEmptyPhones  = errors.New("empty phones")
)
