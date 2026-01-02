package sms

import (
	"context"
	"testing"
)

func TestClient_Send(t *testing.T) {
	// Skip в CI или если нет реальных credentials
	t.Skip("Skip integration test - requires real credentials")

	client := NewClient("", "")

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

func TestStatus_IsError(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"not found", StatusNotFound, true},
		{"stopped", StatusStopped, true},
		{"delivered", StatusDelivered, false},
		{"passed to operator", StatusPassedToOperator, false},
		{"expired", StatusExpired, true},
		{"invalid number", StatusInvalidNumber, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsError(); got != tt.want {
				t.Errorf("Status.IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_IsSuccess(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"delivered", StatusDelivered, true},
		{"read", StatusRead, true},
		{"clicked", StatusClicked, true},
		{"pending", StatusPending, false},
		{"not found", StatusNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsSuccess(); got != tt.want {
				t.Errorf("Status.IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
