package app

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/appenv"
	authC "github.com/Rasikrr/bugsy_backend_monolith/internal/cache/auth"
	smsC "github.com/Rasikrr/bugsy_backend_monolith/internal/cache/sms"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	authS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/core/telegram"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
)

type App struct {
	application.App

	smsCache  smsC.Cache
	authCache authC.Cache

	smsClient sms.Client
	tgClient  telegram.Client

	authService authS.Service
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initCache,
		app.initClients,
		app.initServices,
		app.initHTTP,
	} {
		if err := initFn(ctx); err != nil {
			log.Fatal(ctx, "app init", log.Err(err))
		}
	}
	log.Infof(ctx, "env: %s", app.Config().Environment)
	return app
}

func (a *App) initHTTP(_ context.Context) error {
	http.NewServer(
		a.HTTPServer(),
		a.authService,
	)
	return nil
}

func (a *App) initCache(_ context.Context) error {
	smsSpamTTL, err := a.Config().Variables.GetDuration(appenv.SMSSpamTTL)
	if err != nil {
		return err
	}
	a.smsCache = smsC.NewCache(a.Redis(), smsSpamTTL)

	authCodeTTL, err := a.Config().Variables.GetDuration(appenv.AuthCodeTTL)
	if err != nil {
		return err
	}

	a.authCache = authC.NewCache(a.Redis(), authCodeTTL)
	return nil
}

func (a *App) initServices(_ context.Context) error {
	devSMSChatID, err := a.Config().Variables.GetInt(appenv.DevSMSChatID)
	if err != nil {
		return err
	}
	a.authService = authS.NewService(
		a.Config().Environment,
		a.smsClient,
		a.tgClient,
		a.authCache,
		int64(devSMSChatID),
	)
	return nil
}

func (a *App) initClients(_ context.Context) error {
	token, err := a.Config().Variables.GetString(appenv.DevSMSBotToken)
	if err != nil {
		return err
	}
	a.tgClient, err = telegram.NewTelegramClient(token)
	if err != nil {
		return err
	}

	smsClientLogin, err := a.Config().Variables.GetString(appenv.SMSClientLogin)
	if err != nil {
		return err
	}
	smsClientPassword, err := a.Config().Variables.GetString(appenv.SMSClientPassword)
	if err != nil {
		return err
	}

	a.smsClient = sms.NewClient(smsClientLogin, smsClientPassword, a.smsCache)

	return nil
}
