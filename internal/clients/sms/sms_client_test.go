package sms

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	phones  = []string{"77715275251"}
	message = "мам вызывай такси"
)

func TestSMSClient(t *testing.T) {
	//t.Skip()
	ctx := context.Background()
	godotenv.Load()
	host, login, password := os.Getenv("HOST"), os.Getenv("LOGIN"), os.Getenv("PASSWORD")
	cli := NewClient(host, login, password)

	err := cli.Send(ctx, message, phones...)

	require.NoError(t, err)
}
