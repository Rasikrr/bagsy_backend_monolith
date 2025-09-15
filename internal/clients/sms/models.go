// nolint: unused
package sms

import (
	"fmt"
)

type request struct {
	Login          string         `json:"login"`
	Password       string         `json:"psw"`
	Phones         string         `json:"phones"`
	Message        string         `json:"mes"`
	ResponseFormat responseFormat `json:"fmt"`
}

func createRequestBody(login, password, message string, phone string) request {
	req := request{
		Login:          login,
		Password:       password,
		Message:        message,
		Phones:         phone,
		ResponseFormat: responseFormatJSON,
	}
	return req
}

type response struct {
	Error     string    `json:"error,omitempty"`
	ErrorCode errCodes  `json:"error_code,omitempty"`
	MessageID int       `json:"id"`
	Cnt       int       `json:"cnt"`
	Status    smsStatus `json:"status,omitempty"`
}

func (r response) getError() error {
	if r.Error != "" {
		return fmt.Errorf("error while sending sms: %s | error_code: %d", r.Error, r.ErrorCode)
	}
	return nil
}
