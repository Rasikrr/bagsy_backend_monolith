package whatsapp

import (
	"context"
	"os"
	"testing"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIURL     = "https://7105.api.greenapi.com"
	testMediaURL   = "https://7105.media.greenapi.com"
	testInstanceID = "7105335616"
	testAPIToken   = "4caa49e67772405eaeb99ec4f97124bdf605d90e7e53409489"
)

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)
	err := client.Reboot(ctx)
	require.NoError(t, err)
	err = client.SendMessage(ctx, "77715275251", "test message")
	require.NoError(t, err)
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{
			name:     "phone without plus",
			phone:    "77001234567",
			expected: "77001234567@c.us",
		},
		{
			name:     "phone with plus",
			phone:    "+77001234567",
			expected: "77001234567@c.us",
		},
		{
			name:     "empty phone",
			phone:    "",
			expected: "@c.us",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPhoneNumber(tt.phone)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClient_SendMessage_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	tests := []struct {
		name        string
		phoneNumber string
		message     string
		expectedErr error
	}{
		{
			name:        "empty phone number",
			phoneNumber: "",
			message:     "Test message",
			expectedErr: domainErr.ErrWhatsAppEmptyPhone,
		},
		{
			name:        "empty message",
			phoneNumber: "77001234567",
			message:     "",
			expectedErr: domainErr.ErrWhatsAppEmptyMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SendMessage(context.Background(), tt.phoneNumber, tt.message)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestClient_SendFileByURL_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	tests := []struct {
		name        string
		phoneNumber string
		fileURL     string
		expectedErr error
	}{
		{
			name:        "empty phone number",
			phoneNumber: "",
			fileURL:     "https://example.com/file.jpg",
			expectedErr: domainErr.ErrWhatsAppEmptyPhone,
		},
		{
			name:        "empty file URL",
			phoneNumber: "77001234567",
			fileURL:     "",
			expectedErr: domainErr.ErrWhatsAppEmptyFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SendFileByURL(context.Background(), tt.phoneNumber, tt.fileURL, "")
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestClient_SendFileByUpload_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	tests := []struct {
		name        string
		phoneNumber string
		filePath    string
		expectedErr error
	}{
		{
			name:        "empty phone number",
			phoneNumber: "",
			filePath:    "/path/to/file.jpg",
			expectedErr: domainErr.ErrWhatsAppEmptyPhone,
		},
		{
			name:        "empty file path",
			phoneNumber: "77001234567",
			filePath:    "",
			expectedErr: domainErr.ErrWhatsAppEmptyFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SendFileByUpload(context.Background(), tt.phoneNumber, tt.filePath, "")
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestClient_SendLocation_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	err := client.SendLocation(context.Background(), "", 51.5074, -0.1278, "London")
	require.Error(t, err)
	assert.ErrorIs(t, err, domainErr.ErrWhatsAppEmptyPhone)
}

func TestClient_SendContact_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	tests := []struct {
		name        string
		phoneNumber string
		contact     Contact
		wantErr     bool
	}{
		{
			name:        "empty phone number",
			phoneNumber: "",
			contact: Contact{
				PhoneContact: 77001234567,
				FirstName:    "John",
			},
			wantErr: true,
		},
		{
			name:        "empty contact phone",
			phoneNumber: "77001234567",
			contact: Contact{
				PhoneContact: 0,
				FirstName:    "John",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SendContact(context.Background(), tt.phoneNumber, tt.contact)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClient_DownloadFile_Validation(t *testing.T) {
	client := NewClient(testAPIURL, testMediaURL, testInstanceID, testAPIToken)

	tests := []struct {
		name        string
		chatID      string
		messageID   string
		expectedErr string
	}{
		{
			name:        "empty chat ID",
			chatID:      "",
			messageID:   "msg123",
			expectedErr: "chat ID is required",
		},
		{
			name:        "empty message ID",
			chatID:      "77001234567@c.us",
			messageID:   "",
			expectedErr: "message ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.DownloadFile(context.Background(), tt.chatID, tt.messageID)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// Integration tests (требуют реальных credentials)

func TestClient_SendMessage_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")
	testPhone := os.Getenv("WHATSAPP_TEST_PHONE")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" || testPhone == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	err := client.SendMessage(context.Background(), testPhone, "Test message from automated test")
	require.NoError(t, err)
}

func TestClient_SendFileByURL_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")
	testPhone := os.Getenv("WHATSAPP_TEST_PHONE")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" || testPhone == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	// Example image URL
	imageURL := "https://avatars.githubusercontent.com/u/1?v=4"
	err := client.SendFileByURL(context.Background(), testPhone, imageURL, "Test image caption")
	require.NoError(t, err)
}

func TestClient_SendLocation_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")
	testPhone := os.Getenv("WHATSAPP_TEST_PHONE")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" || testPhone == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	// London coordinates
	err := client.SendLocation(context.Background(), testPhone, 51.5074, -0.1278, "London")
	require.NoError(t, err)
}

func TestClient_SendContact_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")
	testPhone := os.Getenv("WHATSAPP_TEST_PHONE")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" || testPhone == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	contact := Contact{
		PhoneContact: 77001234567,
		FirstName:    "John",
		LastName:     "Doe",
		Company:      "Test Company",
	}

	err := client.SendContact(context.Background(), testPhone, contact)
	require.NoError(t, err)
}

func TestClient_GetStateInstance_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	state, err := client.GetStateInstance(context.Background())
	require.NoError(t, err)
	require.NotNil(t, state)
	assert.NotEmpty(t, state.StateInstance)
	t.Logf("Instance state: %s", state.StateInstance)
}

func TestClient_GetSettings_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	settings, err := client.GetSettings(context.Background())
	require.NoError(t, err)
	require.NotNil(t, settings)
	t.Logf("Settings: %+v", settings)
}

func TestClient_SetSettings_Integration(t *testing.T) {
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_API_ID_INSTANCE")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	settings := Settings{
		IncomingWebhook: true,
		OutgoingWebhook: true,
	}

	err := client.SetSettings(context.Background(), settings)
	require.NoError(t, err)
}
