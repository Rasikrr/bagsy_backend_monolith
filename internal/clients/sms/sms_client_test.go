package sms

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	testPhone = "77068138354"
	message   = "сырники"
)

func TestSMSClient(t *testing.T) {
	//t.Skip()
	ctx := context.Background()
	godotenv.Load()
	login, password := os.Getenv("LOGIN"), os.Getenv("PASSWORD")
	cli := NewClient(login, password, nil)

	err := cli.Send(ctx, testPhone, message)

	require.NoError(t, err)
}
