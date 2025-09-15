package app

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/appenv"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/telegram"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
)

type App struct {
	application.App

	tgDevAuthCache auth.Cache

	tgDevAuthClient telegram.Client
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initHTTP,
		app.initCache,
		app.initClients,
	} {
		if err := initFn(ctx); err != nil {
			log.Fatal(ctx, "app init", log.Err(err))
		}
	}
	return app
}

func (a *App) initHTTP(ctx context.Context) error {
	http.NewServer(a.HTTPServer())
	return nil
}

func (a *App) initCache(ctx context.Context) error {
	authCodeTTL, err := a.Config().Variables.GetDuration(appenv.AuthCodeTTL)
	if err != nil {
		return err
	}
	a.tgDevAuthCache = auth.NewCache(a.Redis(), authCodeTTL)

	return nil
}

func (a *App) initClients(ctx context.Context) error {
	token, err := a.Config().Variables.GetString(appenv.DevSMSBotToken)
	if err != nil {
		return err
	}

	chatID, err := a.Config().Variables.GetInt(appenv.DevSMSChatID)
	if err != nil {
		return err
	}

	a.tgDevAuthClient = telegram.NewClient(
		token,
		int64(chatID),
		"dev_sms_bot",
	)

	return nil
}
