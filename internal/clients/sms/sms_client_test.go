package sms

import (
	"context"
	"os"
	"testing"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/cache/sms"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

var (
	testPhone = "77715275251"
	message   = "сырники"
)

func TestSMSClient(t *testing.T) {
	t.Skip("to save money:)")

	godotenv.Load()
	login, password := os.Getenv("LOGIN"), os.Getenv("PASSWORD")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := sms.NewMockCache(ctrl)
	mockCache.EXPECT().
		IsSpam(gomock.Any(), testPhone, message).Return(false, nil).
		Times(1)

	mockCache.EXPECT().
		Set(gomock.Any(), testPhone, message).Return(nil).
		Times(1)

	cli := NewClient(login, password, mockCache)
	err := cli.Send(context.Background(), testPhone, message)

	require.NoError(t, err)
}
