package whatsapp

import (
	"context"
	"errors"
	"fmt"

	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)

type Client interface {
	SendMessage(ctx context.Context, phoneNumber, message string) error
}

type client struct {
	api *greenapi.GreenAPI
}

func NewClient(apiURL, mediaURL, idInstance, apiToken string) Client {
	api := &greenapi.GreenAPI{
		APIURL:           apiURL,
		MediaURL:         mediaURL,
		IDInstance:       idInstance,
		APITokenInstance: apiToken,
	}

	return &client{
		api: api,
	}
}

func (c *client) SendMessage(_ context.Context, phoneNumber, message string) error {
	if phoneNumber == "" {
		return errors.New("phone number is required")
	}
	if message == "" {
		return errors.New("message is required")
	}

	// Форматируем номер телефона в формат WhatsApp (например: 11001234567@c.us)
	chatID := formatPhoneNumber(phoneNumber)

	resp, err := c.api.Sending().SendMessage(chatID, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if resp == nil {
		return errors.New("empty response from API")
	}

	return nil
}

// formatPhoneNumber форматирует номер телефона в формат WhatsApp.
// Ожидается номер в международном формате без +, например: 77001234567.
func formatPhoneNumber(phone string) string {
	// Убираем + если есть
	if len(phone) > 0 && phone[0] == '+' {
		phone = phone[1:]
	}
	return phone + "@c.us"
}
