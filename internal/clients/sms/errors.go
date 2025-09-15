package sms

import "errors"

var (
	errEmptyMessage = errors.New("empty message")
	errEmptyPhones  = errors.New("empty phones")
	errCheckStatus  = errors.New("invalid status")
	errSpam         = errors.New("spam, try later")
)
