package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"
	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)

type Client struct {
	api *greenapi.GreenAPI
}

func NewClient(apiURL, mediaURL, idInstance, apiToken string) *Client {
	api := &greenapi.GreenAPI{
		APIURL:           apiURL,
		MediaURL:         mediaURL,
		IDInstance:       idInstance,
		APITokenInstance: apiToken,
	}

	return &Client{
		api: api,
	}
}

func (c *Client) SendMessage(_ context.Context, phoneNumber, message string) error {
	if phoneNumber == "" {
		return ErrEmptyPhone
	}
	if message == "" {
		return ErrEmptyMessage
	}

	chatID := formatPhoneNumber(phoneNumber)

	resp, err := c.api.Sending().SendMessage(chatID, message)
	if err != nil {
		return errors.Wrap(ErrSendFailed, err.Error())
	}

	if resp == nil {
		return ErrEmptyResponse
	}

	return nil
}

func (c *Client) SendFileByURL(_ context.Context, phoneNumber, fileURL, caption string) error {
	if phoneNumber == "" {
		return ErrEmptyPhone
	}
	if fileURL == "" {
		return ErrEmptyFile
	}

	chatID := formatPhoneNumber(phoneNumber)

	var opts []greenapi.SendFileByUrlOption
	if caption != "" {
		opts = append(opts, greenapi.OptionalCaptionSendUrl(caption))
	}

	resp, err := c.api.Sending().SendFileByUrl(chatID, fileURL, "", opts...)
	if err != nil {
		return errors.Wrap(ErrSendFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return ErrEmptyResponse
	}

	return nil
}

func (c *Client) SendFileByUpload(_ context.Context, phoneNumber, filePath, caption string) error {
	if phoneNumber == "" {
		return ErrEmptyPhone
	}
	if filePath == "" {
		return ErrEmptyFile
	}

	chatID := formatPhoneNumber(phoneNumber)

	var opts []greenapi.SendFileByUploadOption
	if caption != "" {
		opts = append(opts, greenapi.OptionalCaptionSendUpload(caption))
	}

	resp, err := c.api.Sending().SendFileByUpload(chatID, filePath, "", opts...)
	if err != nil {
		return errors.Wrap(ErrSendFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return ErrEmptyResponse
	}

	return nil
}

func (c *Client) SendLocation(_ context.Context, phoneNumber string, latitude, longitude float64, locationName string) error {
	if phoneNumber == "" {
		return ErrEmptyPhone
	}

	chatID := formatPhoneNumber(phoneNumber)

	var opts []greenapi.SendLocationOption
	if locationName != "" {
		opts = append(opts, greenapi.OptionalNameLocation(locationName))
	}

	resp, err := c.api.Sending().SendLocation(chatID, float32(latitude), float32(longitude), opts...)
	if err != nil {
		return errors.Wrap(ErrSendFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return ErrEmptyResponse
	}

	return nil
}

func (c *Client) SendContact(_ context.Context, phoneNumber string, contact Contact) error {
	if phoneNumber == "" {
		return ErrEmptyPhone
	}
	if contact.PhoneContact == 0 {
		return ErrEmptyContact
	}

	chatID := formatPhoneNumber(phoneNumber)

	greenContact := greenapi.Contact{
		PhoneContact: contact.PhoneContact,
		FirstName:    contact.FirstName,
		LastName:     contact.LastName,
		MiddleName:   contact.MiddleName,
		Company:      contact.Company,
	}

	resp, err := c.api.Sending().SendContact(chatID, greenContact)
	if err != nil {
		return errors.Wrap(ErrSendFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return ErrEmptyResponse
	}

	return nil
}

// GetStateInstance получает состояние инстанса
func (c *Client) GetStateInstance(_ context.Context) (*StateInstance, error) {
	resp, err := c.api.Account().GetStateInstance()
	if err != nil {
		return nil, errors.Wrap(ErrGetStateFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, ErrEmptyResponse
	}

	var state StateInstance
	if unmarshalErr := json.Unmarshal(resp.Body, &state); unmarshalErr != nil {
		return nil, errors.Wrap(ErrUnmarshalFailed, unmarshalErr.Error())
	}

	return &state, nil
}

// GetSettings получает настройки инстанса
func (c *Client) GetSettings(_ context.Context) (*Settings, error) {
	resp, err := c.api.Account().GetSettings()
	if err != nil {
		return nil, errors.Wrap(ErrGetSettingsFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, ErrEmptyResponse
	}

	var settings Settings
	if unmarshalErr := json.Unmarshal(resp.Body, &settings); unmarshalErr != nil {
		return nil, errors.Wrap(ErrUnmarshalFailed, unmarshalErr.Error())
	}

	return &settings, nil
}

// SetSettings устанавливает настройки инстанса
func (c *Client) SetSettings(_ context.Context, settings Settings) error {
	var opts []greenapi.SetSettingsOption

	if settings.WebhookURL != "" {
		opts = append(opts, greenapi.OptionalWebhookUrl(settings.WebhookURL))
	}
	if settings.WebhookURLToken != "" {
		opts = append(opts, greenapi.OptionalWebhookUrlToken(settings.WebhookURLToken))
	}
	if settings.DelaySendMessagesMS > 0 {
		opts = append(opts, greenapi.OptionalDelaySendMessages(uint(settings.DelaySendMessagesMS)))
	}
	if settings.MarkIncomingMsgReaded {
		opts = append(opts, greenapi.OptionalMarkIncomingMessagesRead(settings.MarkIncomingMsgReaded))
	}
	if settings.OutgoingWebhook {
		opts = append(opts, greenapi.OptionalOutgoingWebhook(settings.OutgoingWebhook))
	}
	if settings.IncomingWebhook {
		opts = append(opts, greenapi.OptionalIncomingWebhook(settings.IncomingWebhook))
	}

	resp, err := c.api.Account().SetSettings(opts...)
	if err != nil {
		return errors.Wrap(ErrSetSettingsFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return ErrEmptyResponse
	}

	return nil
}

// Reboot перезагружает инстанс
func (c *Client) Reboot(_ context.Context) error {
	resp, err := c.api.Account().Reboot()
	if err != nil {
		return errors.Wrap(ErrRebootFailed, err.Error())
	}
	if resp == nil {
		return ErrEmptyResponse
	}

	return nil
}

// Logout выполняет выход из аккаунта
func (c *Client) Logout(_ context.Context) error {
	resp, err := c.api.Account().Logout()
	if err != nil {
		return errors.Wrap(ErrLogoutFailed, err.Error())
	}

	if resp == nil {
		return ErrEmptyResponse
	}

	return nil
}

// DownloadFile загружает файл из сообщения
func (c *Client) DownloadFile(_ context.Context, chatID, messageID string) ([]byte, error) {
	if chatID == "" {
		return nil, ErrEmptyChatID
	}
	if messageID == "" {
		return nil, ErrEmptyMsgID
	}

	resp, err := c.api.Receiving().DownloadFile(chatID, messageID)
	if err != nil {
		return nil, errors.Wrap(ErrDownloadFailed, err.Error())
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, ErrEmptyResponse
	}

	return resp.Body, nil
}

func formatPhoneNumber(phone string) string {
	if len(phone) > 0 && phone[0] == '+' {
		phone = phone[1:]
	}
	return phone + "@c.us"
}
