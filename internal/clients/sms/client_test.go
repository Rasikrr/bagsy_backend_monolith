package sms

import (
	"context"
	"testing"
)

func TestClient_Send(t *testing.T) {
	// Skip в CI или если нет реальных credentials
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := NewClient("dbagsy25", "bJwwe5@B!mdZ3Bz")

	tests := []struct {
		name    string
		phone   string
		message string
		wantErr bool
	}{
		{
			name:    "empty message",
			phone:   "79999999999",
			message: "",
			wantErr: true,
		},
		{
			name:    "empty phone",
			phone:   "",
			message: "test message",
			wantErr: true,
		},
		{
			name:    "valid request (will fail without real credentials)",
			phone:   "77715275251",
			message: "Test message",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Send(context.Background(), tt.phone, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSMSStatus_IsError(t *testing.T) {
	tests := []struct {
		name   string
		status SMSStatus
		want   bool
	}{
		{"not found", SMSStatusNotFound, true},
		{"stopped", SMSStatusStopped, true},
		{"delivered", SMSStatusDelivered, false},
		{"passed to operator", SMSStatusPassedToOperator, false},
		{"expired", SMSStatusExpired, true},
		{"invalid number", SMSStatusInvalidNumber, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsError(); got != tt.want {
				t.Errorf("SMSStatus.IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSMSStatus_IsSuccess(t *testing.T) {
	tests := []struct {
		name   string
		status SMSStatus
		want   bool
	}{
		{"delivered", SMSStatusDelivered, true},
		{"read", SMSStatusRead, true},
		{"clicked", SMSStatusClicked, true},
		{"pending", SMSStatusPending, false},
		{"not found", SMSStatusNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsSuccess(); got != tt.want {
				t.Errorf("SMSStatus.IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
