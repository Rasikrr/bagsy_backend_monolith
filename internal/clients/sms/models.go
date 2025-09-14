package sms

import "strings"

type request struct {
	Login    string `json:"login"`
	Password string `json:"psw"`
	Phones   string `json:"phones"`
	Message  string `json:"mes"`
}

func createRequestBody(login, password, message string, phones []string) request {
	req := request{
		Login:    login,
		Password: password,
		Message:  message,
		Phones:   strings.Join(phones, ","),
	}
	return req
}
