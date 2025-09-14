// nolint: unused
package sms

import (
	"fmt"
	"strings"
)

type responseFormat uint8

const (
	responseFormatString responseFormat = iota
	responseFormatNumber
	responseFormatXML
	responseFormatJSON
)

type smsStatus int8

const (
	smsStatusNotFound smsStatus = iota - 3
	smsStatusStopped
	smsStatusPending
	smsStatusPassedToOperator
	smsSatusDelivered
	smsStatusChecked
)

type request struct {
	Login          string         `json:"login"`
	Password       string         `json:"psw"`
	Phones         string         `json:"phones"`
	Message        string         `json:"mes"`
	ResponseFormat responseFormat `json:"fmt"`
}

func createRequestBody(login, password, message string, phones []string) request {
	req := request{
		Login:          login,
		Password:       password,
		Message:        message,
		Phones:         strings.Join(phones, ","),
		ResponseFormat: responseFormatJSON,
	}
	return req
}

type sendResponse struct {
	Error     string    `json:"error,omitempty"`
	ErrorCode int       `json:"error_code,omitempty"`
	MessageID int       `json:"id"`
	Cnt       int       `json:"cnt"`
	Status    smsStatus `json:"status,omitempty"`
}

func (s sendResponse) getError() error {
	if s.Error != "" {
		return fmt.Errorf("error while sending sms: %s | error_code: %d", s.Error, s.ErrorCode)
	}
	return nil
}
