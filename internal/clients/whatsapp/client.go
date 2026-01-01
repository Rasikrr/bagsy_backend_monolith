package whatsapp

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//
//	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
//)
//
//type Client struct {
//	api *greenapi.GreenAPI
//}
//
//func NewClient(apiURL, mediaURL, idInstance, apiToken string) *Client {
//	api := &greenapi.GreenAPI{
//		APIURL:           apiURL,
//		MediaURL:         mediaURL,
//		IDInstance:       idInstance,
//		APITokenInstance: apiToken,
//	}
//
//	return &Client{
//		api: api,
//	}
//}
//
//func (c *Client) SendMessage(_ context.Context, phoneNumber, message string) error {
//	if phoneNumber == "" {
//		return ErrWhatsAppPhoneRequired
//	}
//	if message == "" {
//		return ErrWhatsAppMessageRequired
//	}
//
//	chatID := formatPhoneNumber(phoneNumber)
//
//	resp, err := c.api.Sending().SendMessage(chatID, message)
//	if err != nil {
//		return fmt.Errorf("%w: %w", domainErrors.ErrWhatsAppSendFailed, err)
//	}
//
//	if resp == nil {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//func (c *Client) SendFileByURL(_ context.Context, phoneNumber, fileURL, caption string) error {
//	if phoneNumber == "" {
//		return domainErrors.ErrWhatsAppPhoneRequired
//	}
//	if fileURL == "" {
//		return domainErrors.ErrWhatsAppFileRequired
//	}
//
//	chatID := formatPhoneNumber(phoneNumber)
//
//	var opts []greenapi.SendFileByUrlOption
//	if caption != "" {
//		opts = append(opts, greenapi.OptionalCaptionSendUrl(caption))
//	}
//
//	resp, err := c.api.Sending().SendFileByUrl(chatID, fileURL, "", opts...)
//	if err != nil {
//		return fmt.Errorf("%w: %w", domainErrors.ErrWhatsAppSendFailed, err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//func (c *Client) SendFileByUpload(_ context.Context, phoneNumber, filePath, caption string) error {
//	if phoneNumber == "" {
//		return domainErrors.ErrWhatsAppPhoneRequired
//	}
//	if filePath == "" {
//		return domainErrors.ErrWhatsAppFileRequired
//	}
//
//	chatID := formatPhoneNumber(phoneNumber)
//
//	var opts []greenapi.SendFileByUploadOption
//	if caption != "" {
//		opts = append(opts, greenapi.OptionalCaptionSendUpload(caption))
//	}
//
//	resp, err := c.api.Sending().SendFileByUpload(chatID, filePath, "", opts...)
//	if err != nil {
//		return fmt.Errorf("%w: %w", domainErrors.ErrWhatsAppSendFailed, err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//func (c *Client) SendLocation(_ context.Context, phoneNumber string, latitude, longitude float64, locationName string) error {
//	if phoneNumber == "" {
//		return domainErrors.ErrWhatsAppPhoneRequired
//	}
//
//	chatID := formatPhoneNumber(phoneNumber)
//
//	var opts []greenapi.SendLocationOption
//	if locationName != "" {
//		opts = append(opts, greenapi.OptionalNameLocation(locationName))
//	}
//
//	resp, err := c.api.Sending().SendLocation(chatID, float32(latitude), float32(longitude), opts...)
//	if err != nil {
//		return fmt.Errorf("%w: %w", domainErrors.ErrWhatsAppSendFailed, err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//func (c *Client) SendContact(_ context.Context, phoneNumber string, contact Contact) error {
//	if phoneNumber == "" {
//		return domainErrors.ErrWhatsAppPhoneRequired
//	}
//	if contact.PhoneContact == 0 {
//		return fmt.Errorf("contact phone number is required")
//	}
//
//	chatID := formatPhoneNumber(phoneNumber)
//
//	greenContact := greenapi.Contact{
//		PhoneContact: contact.PhoneContact,
//		FirstName:    contact.FirstName,
//		LastName:     contact.LastName,
//		MiddleName:   contact.MiddleName,
//		Company:      contact.Company,
//	}
//
//	resp, err := c.api.Sending().SendContact(chatID, greenContact)
//	if err != nil {
//		return fmt.Errorf("%w: %w", domainErrors.ErrWhatsAppSendFailed, err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//// GetStateInstance получает состояние инстанса
//func (c *Client) GetStateInstance(_ context.Context) (*StateInstance, error) {
//	resp, err := c.api.Account().GetStateInstance()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get state instance: %w", err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return nil, domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	var state StateInstance
//	if err := json.Unmarshal(resp.Body, &state); err != nil {
//		return nil, fmt.Errorf("failed to unmarshal state response: %w", err)
//	}
//
//	return &state, nil
//}
//
//// GetSettings получает настройки инстанса
//func (c *Client) GetSettings(_ context.Context) (*Settings, error) {
//	resp, err := c.api.Account().GetSettings()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get settings: %w", err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return nil, domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	var settings Settings
//	if err := json.Unmarshal(resp.Body, &settings); err != nil {
//		return nil, fmt.Errorf("failed to unmarshal settings response: %w", err)
//	}
//
//	return &settings, nil
//}
//
//// SetSettings устанавливает настройки инстанса
//func (c *Client) SetSettings(_ context.Context, settings Settings) error {
//	var opts []greenapi.SetSettingsOption
//
//	if settings.WebhookURL != "" {
//		opts = append(opts, greenapi.OptionalWebhookUrl(settings.WebhookURL))
//	}
//	if settings.WebhookURLToken != "" {
//		opts = append(opts, greenapi.OptionalWebhookUrlToken(settings.WebhookURLToken))
//	}
//	if settings.DelaySendMessagesMS > 0 {
//		opts = append(opts, greenapi.OptionalDelaySendMessages(uint(settings.DelaySendMessagesMS)))
//	}
//	if settings.MarkIncomingMsgReaded {
//		opts = append(opts, greenapi.OptionalMarkIncomingMessagesRead(settings.MarkIncomingMsgReaded))
//	}
//	if settings.OutgoingWebhook {
//		opts = append(opts, greenapi.OptionalOutgoingWebhook(settings.OutgoingWebhook))
//	}
//	if settings.IncomingWebhook {
//		opts = append(opts, greenapi.OptionalIncomingWebhook(settings.IncomingWebhook))
//	}
//
//	resp, err := c.api.Account().SetSettings(opts...)
//	if err != nil {
//		return fmt.Errorf("failed to set settings: %w", err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//// Reboot перезагружает инстанс
//func (c *Client) Reboot(_ context.Context) error {
//	resp, err := c.api.Account().Reboot()
//	if err != nil {
//		return fmt.Errorf("failed to reboot instance: %w", err)
//	}
//
//	if resp == nil {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//// Logout выполняет выход из аккаунта
//func (c *Client) Logout(_ context.Context) error {
//	resp, err := c.api.Account().Logout()
//	if err != nil {
//		return fmt.Errorf("failed to logout: %w", err)
//	}
//
//	if resp == nil {
//		return domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return nil
//}
//
//// DownloadFile загружает файл из сообщения
//func (c *Client) DownloadFile(_ context.Context, chatID, messageID string) ([]byte, error) {
//	if chatID == "" {
//		return nil, fmt.Errorf("chat ID is required")
//	}
//	if messageID == "" {
//		return nil, fmt.Errorf("message ID is required")
//	}
//
//	resp, err := c.api.Receiving().DownloadFile(chatID, messageID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to download file: %w", err)
//	}
//
//	if resp == nil || resp.StatusCode != 200 {
//		return nil, domainErrors.ErrWhatsAppEmptyResponse
//	}
//
//	return resp.Body, nil
//}
//
//func formatPhoneNumber(phone string) string {
//	if len(phone) > 0 && phone[0] == '+' {
//		phone = phone[1:]
//	}
//	return phone + "@c.us"
//}
