package telegram

import (
	"context"
	"github.com/Rasikrr/core/log"
	"time"

	"github.com/avast/retry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client interface {
	Send(ctx context.Context, msg string) error
}

type client struct {
	botName string
	chatID  int64
	token   string
}

func NewClient(token string, chatID int64, botName string) Client {
	return &client{
		botName: botName,
		chatID:  chatID,
		token:   token,
	}
}

func (s *client) Send(ctx context.Context, msg string) error {
	bot, err := tgbotapi.NewBotAPI(s.token)
	if err != nil {
		return err
	}

	message := tgbotapi.NewMessage(s.chatID, msg)
	message.ParseMode = tgbotapi.ModeMarkdownV2

	return retry.Do(func() error {
		_, err = bot.Send(message)
		if err != nil {
			log.Error(ctx, "error while sending timer", log.Err(err))
			return err
		}
		return err
	}, retry.Context(ctx), retry.Attempts(3), retry.Delay(10*time.Second),
	)
}
