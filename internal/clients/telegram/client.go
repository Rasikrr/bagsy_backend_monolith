package telegram

import (
	"context"
	"strconv"

	"github.com/Rasikrr/core/telegram"
)

type Client struct {
	cli telegram.Client
}

func NewClient(cli telegram.Client) *Client {
	return &Client{
		cli: cli,
	}
}

func (c *Client) Send(ctx context.Context, chatID, token string) error {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}
	return c.cli.SendText(ctx, chatIDInt, token)
}
