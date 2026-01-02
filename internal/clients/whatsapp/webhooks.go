package whatsapp

import "time"

// WebhookNotification представляет базовую структуру webhook уведомления
type WebhookNotification struct {
	TypeWebhook       string        `json:"typeWebhook"`
	InstanceData      Instance      `json:"instanceData"`
	Timestamp         int64         `json:"timestamp"`
	IDMessage         string        `json:"idMessage"`
	SenderData        SenderData    `json:"senderData,omitempty"`
	MessageData       MessageData   `json:"messageData,omitempty"`
	StatusData        StatusData    `json:"statusData,omitempty"`
	StateInstanceData StateInstance `json:"stateInstanceData,omitempty"`
}

// Instance представляет информацию об инстансе
type Instance struct {
	IDInstance   string `json:"idInstance"`
	WID          string `json:"wid"`
	TypeInstance string `json:"typeInstance"`
}

// SenderData представляет информацию об отправителе
type SenderData struct {
	ChatID            string `json:"chatId"`
	ChatName          string `json:"chatName"`
	Sender            string `json:"sender"`
	SenderName        string `json:"senderName"`
	SenderContactName string `json:"senderContactName"`
}

// MessageData представляет данные сообщения
type MessageData struct {
	TypeMessage             string                   `json:"typeMessage"`
	TextMessageData         *TextMessageData         `json:"textMessageData,omitempty"`
	ExtendedTextMessageData *ExtendedTextMessageData `json:"extendedTextMessageData,omitempty"`
	ImageMessageData        *FileMessageData         `json:"imageMessageData,omitempty"`
	VideoMessageData        *FileMessageData         `json:"videoMessageData,omitempty"`
	DocumentMessageData     *FileMessageData         `json:"documentMessageData,omitempty"`
	AudioMessageData        *FileMessageData         `json:"audioMessageData,omitempty"`
	LocationMessageData     *LocationMessageData     `json:"locationMessageData,omitempty"`
	ContactMessageData      *ContactMessageData      `json:"contactMessageData,omitempty"`
	QuotedMessage           *QuotedMessage           `json:"quotedMessage,omitempty"`
	IsForwarded             bool                     `json:"isForwarded"`
}

// TextMessageData представляет текстовое сообщение
type TextMessageData struct {
	TextMessage string `json:"textMessage"`
}

// ExtendedTextMessageData представляет расширенное текстовое сообщение
type ExtendedTextMessageData struct {
	Text          string `json:"text"`
	Description   string `json:"description"`
	Title         string `json:"title"`
	PreviewType   string `json:"previewType"`
	JPEGThumbnail string `json:"jpegThumbnail"`
	StanzaID      string `json:"stanzaId"`
	Participant   string `json:"participant"`
}

// FileMessageData представляет файловое сообщение (изображение, видео, документ, аудио)
type FileMessageData struct {
	DownloadURL   string `json:"downloadUrl"`
	Caption       string `json:"caption"`
	FileName      string `json:"fileName"`
	JPEGThumbnail string `json:"jpegThumbnail"`
	MimeType      string `json:"mimeType"`
	FileLength    int64  `json:"fileLength"`
}

// LocationMessageData представляет сообщение с геолокацией
type LocationMessageData struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	JPEGThumbnail string  `json:"jpegThumbnail"`
	NameLocation  string  `json:"nameLocation"`
	Address       string  `json:"address"`
}

// ContactMessageData представляет сообщение с контактом
type ContactMessageData struct {
	DisplayName string       `json:"displayName"`
	VCard       string       `json:"vcard"`
	Contact     ContactVCard `json:"contact"`
}

// ContactVCard представляет данные контакта из vCard
type ContactVCard struct {
	DisplayName string `json:"displayName"`
	VCard       string `json:"vcard"`
}

// QuotedMessage представляет цитируемое сообщение
type QuotedMessage struct {
	StanzaID    string `json:"stanzaId"`
	Participant string `json:"participant"`
	TypeMessage string `json:"typeMessage"`
}

// StatusData представляет статус сообщения (доставлено, прочитано и т.д.)
type StatusData struct {
	Status    string `json:"status"` // sent, delivered, read, failed
	Timestamp int64  `json:"timestamp"`
}

// WebhookType константы типов webhook уведомлений
const (
	// Входящие сообщения
	WebhookTypeIncomingMessageReceived = "incomingMessageReceived"
	WebhookTypeIncomingCall            = "incomingCall"

	// Статусы сообщений
	WebhookTypeOutgoingMessageStatus    = "outgoingMessageStatus"
	WebhookTypeOutgoingAPIMessageStatus = "outgoingAPIMessageStatus"

	// Состояние устройства
	WebhookTypeDeviceInfo           = "deviceInfo"
	WebhookTypeStateInstanceChanged = "stateInstanceChanged"

	// Статусы участников
	WebhookTypeStatusInstanceChanged = "statusInstanceChanged"
)

// MessageType константы типов сообщений
const (
	MessageTypeText          = "textMessage"
	MessageTypeExtendedText  = "extendedTextMessage"
	MessageTypeImage         = "imageMessage"
	MessageTypeVideo         = "videoMessage"
	MessageTypeDocument      = "documentMessage"
	MessageTypeAudio         = "audioMessage"
	MessageTypeVoice         = "voiceMessage"
	MessageTypeLocation      = "locationMessage"
	MessageTypeContact       = "contactMessage"
	MessageTypeContactsArray = "contactsArrayMessage"
	MessageTypeSticker       = "stickerMessage"
	MessageTypePoll          = "pollMessage"
)

// MessageStatus константы статусов сообщений
const (
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
	MessageStatusFailed    = "failed"
	MessageStatusDeleted   = "deleted"
)

// StateInstanceStatus константы состояний инстанса
const (
	StateNotAuthorized = "notAuthorized"
	StateAuthorized    = "authorized"
	StateBlocked       = "blocked"
	StateSleepMode     = "sleepMode"
	StateStarting      = "starting"
	StateYellowCard    = "yellowCard"
)

// GetMessageTime возвращает время сообщения
func (w *WebhookNotification) GetMessageTime() time.Time {
	return time.Unix(w.Timestamp, 0)
}

// GetStatusTime возвращает время статуса
func (s *StatusData) GetStatusTime() time.Time {
	return time.Unix(s.Timestamp, 0)
}

// IsTextMessage проверяет, является ли сообщение текстовым
func (m *MessageData) IsTextMessage() bool {
	return m.TypeMessage == MessageTypeText || m.TypeMessage == MessageTypeExtendedText
}

// IsMediaMessage проверяет, является ли сообщение медиа-файлом
func (m *MessageData) IsMediaMessage() bool {
	return m.TypeMessage == MessageTypeImage ||
		m.TypeMessage == MessageTypeVideo ||
		m.TypeMessage == MessageTypeDocument ||
		m.TypeMessage == MessageTypeAudio ||
		m.TypeMessage == MessageTypeVoice
}

// IsLocationMessage проверяет, является ли сообщение геолокацией
func (m *MessageData) IsLocationMessage() bool {
	return m.TypeMessage == MessageTypeLocation
}

// IsContactMessage проверяет, является ли сообщение контактом
func (m *MessageData) IsContactMessage() bool {
	return m.TypeMessage == MessageTypeContact || m.TypeMessage == MessageTypeContactsArray
}

// GetText возвращает текст сообщения (если есть)
func (m *MessageData) GetText() string {
	if m.TextMessageData != nil {
		return m.TextMessageData.TextMessage
	}
	if m.ExtendedTextMessageData != nil {
		return m.ExtendedTextMessageData.Text
	}
	return ""
}

// GetCaption возвращает подпись медиа-файла (если есть)
func (m *MessageData) GetCaption() string {
	if m.ImageMessageData != nil {
		return m.ImageMessageData.Caption
	}
	if m.VideoMessageData != nil {
		return m.VideoMessageData.Caption
	}
	if m.DocumentMessageData != nil {
		return m.DocumentMessageData.Caption
	}
	if m.AudioMessageData != nil {
		return m.AudioMessageData.Caption
	}
	return ""
}

// GetFileURL возвращает URL файла (если это медиа-сообщение)
func (m *MessageData) GetFileURL() string {
	if m.ImageMessageData != nil {
		return m.ImageMessageData.DownloadURL
	}
	if m.VideoMessageData != nil {
		return m.VideoMessageData.DownloadURL
	}
	if m.DocumentMessageData != nil {
		return m.DocumentMessageData.DownloadURL
	}
	if m.AudioMessageData != nil {
		return m.AudioMessageData.DownloadURL
	}
	return ""
}
