package whatsapp

type Contact struct {
	PhoneContact int    `json:"phoneContact"` // Номер телефона контакта
	FirstName    string `json:"firstName"`    // Имя
	LastName     string `json:"lastName"`     // Фамилия
	MiddleName   string `json:"middleName"`   // Отчество
	Company      string `json:"company"`      // Компания
}

type StateInstance struct {
	StateInstance string `json:"stateInstance"` // notAuthorized, authorized, blocked, sleepMode, starting, yellowCard
}

type Settings struct {
	WebhookURL            string `json:"webhookUrl,omitempty"`
	WebhookURLToken       string `json:"webhookUrlToken,omitempty"`
	DelaySendMessagesMS   int    `json:"delaySendMessagesMilliseconds,omitempty"`
	MarkIncomingMsgReaded bool   `json:"markIncomingMessagesReaded,omitempty"`
	ProxyInstance         string `json:"proxyInstance,omitempty"`
	OutgoingWebhook       bool   `json:"outgoingWebhook,omitempty"`
	IncomingWebhook       bool   `json:"incomingWebhook,omitempty"`
}
