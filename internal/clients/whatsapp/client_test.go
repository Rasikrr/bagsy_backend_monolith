package whatsapp

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

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
			if result != tt.expected {
				t.Errorf("formatPhoneNumber(%q) = %q, want %q", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestClient_SendMessage(t *testing.T) {
	// Этот тест требует реальные credentials для работы
	// Раскомментируйте и добавьте свои данные для интеграционного теста
	t.Skip("Skip integration test - requires real credentials")

	_ = godotenv.Load()

	apiURL := os.Getenv("WHATSAPP_API_URL")
	mediaURL := os.Getenv("WHATSAPP_MEDIA_URL")
	instanceID := os.Getenv("WHATSAPP_INSTANCE_ID")
	apiToken := os.Getenv("WHATSAPP_API_TOKEN")
	testPhone := os.Getenv("WHATSAPP_TEST_PHONE")

	if apiURL == "" || mediaURL == "" || instanceID == "" || apiToken == "" || testPhone == "" {
		t.Skip("Skip integration test - missing environment variables")
	}

	client := NewClient(apiURL, mediaURL, instanceID, apiToken)

	err := client.SendMessage(context.Background(), testPhone, "Test message")
	if err != nil {
		t.Errorf("SendMessage() error = %v", err)
	}
}

func TestClient_SendMessage_EmptyPhone(t *testing.T) {
	client := NewClient("https://7105.api.greenapi.com", "https://7105.media.greenapi.com", "test_id", "test_token")

	err := client.SendMessage(context.Background(), "", "Test message")
	if err == nil {
		t.Error("SendMessage() expected error for empty phone, got nil")
	}
}

func TestClient_SendMessage_EmptyMessage(t *testing.T) {
	client := NewClient("https://7105.api.greenapi.com", "https://7105.media.greenapi.com", "test_id", "test_token")

	err := client.SendMessage(context.Background(), "77001234567", "")
	if err == nil {
		t.Error("SendMessage() expected error for empty message, got nil")
	}
}
