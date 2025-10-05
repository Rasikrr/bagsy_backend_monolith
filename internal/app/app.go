package app

import (
	"context"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/appenv"
	authC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/auth"
	smsC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	bagsiesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	usersR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	bagsiesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	formsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"

	"github.com/Rasikrr/core/telegram"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
)

type App struct {
	application.App

	smsCache  smsC.Cache
	authCache authC.Cache

	smsClient      sms.Client
	tgClient       telegram.Client
	whatsAppClient whatsapp.Client

	usersRepo      usersR.Repository
	bagsiesRepo    bagsiesR.Repository
	bagsiesService bagsiesS.Service
	formsRepo      forms.Repository

	authService  authS.Service
	formsService formsS.Service
	usersService usersS.Service
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initCache,
		app.initClients,
		app.initRepositories,
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
	swaggerHost, err := a.Config().Variables.GetString(appenv.SwaggerHost)
	if err != nil {
		return err
	}
	swaggerScheme, err := a.Config().Variables.GetString(appenv.SwaggerScheme)
	if err != nil {
		return err
	}
	http.NewServer(
		a.HTTPServer(),
		swaggerHost,
		swaggerScheme,
		a.authService,
		a.usersService,
		a.formsService,
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

func (a *App) initRepositories(_ context.Context) error {
	a.usersRepo = usersR.NewRepository(a.Postgres())
	return nil
}

func (a *App) initServices(_ context.Context) error {
	devSMSChatID, err := a.Config().Variables.GetInt(appenv.DevSMSChatID)
	if err != nil {
		return err
	}
	authConfirmationURL, err := a.Config().Variables.GetString(appenv.AuthConfirmationURL)
	if err != nil {
		return err
	}
	jwtSecret, err := a.Config().Variables.GetString(appenv.JWTSecret)
	if err != nil {
		return err
	}

	a.usersService = usersS.NewService(a.usersRepo)

	a.authService = authS.NewService(
		a.smsClient,
		a.whatsAppClient,
		a.tgClient,
		a.authCache,
		a.usersService,
		a.PostgresTXManager(),
		int64(devSMSChatID),
		authConfirmationURL,
		jwtSecret,
	)

	a.bagsiesService = bagsiesS.NewService(
		a.bagsiesRepo,
	)

	a.formsService = formsS.NewService(a.formsRepo)
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

	whatsappAPIURL, err := a.Config().Variables.GetString(appenv.WhatsAppAPIURL)
	if err != nil {
		return err
	}

	whatsappMediaURL, err := a.Config().Variables.GetString(appenv.WhatsAppMediaURL)
	if err != nil {
		return err
	}

	whatsappIDInstance, err := a.Config().Variables.GetString(appenv.WhatsAppIDInstance)
	if err != nil {
		return err
	}

	whatsappAPIToken, err := a.Config().Variables.GetString(appenv.WhatsAppAPIToken)
	if err != nil {
		return err
	}

	a.whatsAppClient = whatsapp.NewClient(whatsappAPIURL, whatsappMediaURL, whatsappIDInstance, whatsappAPIToken)
	return nil
}
