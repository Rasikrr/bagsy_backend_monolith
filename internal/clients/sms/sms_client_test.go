package sms

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	testPhone = "77783784148"
	message   = "сырники"
)

func TestSMSClient(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	godotenv.Load()
	login, password := os.Getenv("LOGIN"), os.Getenv("PASSWORD")
	cli := NewClient(login, password)

	err := cli.Send(ctx, message, testPhone)

	require.NoError(t, err)
}
