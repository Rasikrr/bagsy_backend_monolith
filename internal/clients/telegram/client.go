package telegram

import (
	"context"

	"github.com/Rasikrr/core/telegram"
)

type Client struct {
	cli    telegram.Client
	chatID int64
}

func NewClient(cli telegram.Client, chatID int64) *Client {
	return &Client{
		cli:    cli,
		chatID: chatID,
	}
}

func (c *Client) SendMessage(ctx context.Context, message string) error {
	return c.cli.SendText(ctx, c.chatID, message)
}
