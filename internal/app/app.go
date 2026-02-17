package app

import (
	"context"
	"time"

	appenv "github.com/Rasikrr/bagsy_backend_monolith/internal/appenvs"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/s3"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/whatsapp"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
	"github.com/robfig/cron/v3"
)

type App struct {
	application.App

	smsClient      *sms.Client
	whatsappClient *whatsapp.Client
	s3Client       *s3.Client
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initInfra,
		app.initRepositories,
		app.initClients,
		app.initServices,
		app.initHTTP,
		app.initJobs,
	} {
		if err := initFn(ctx); err != nil {
			log.Fatal(ctx, "app init", log.Err(err))
		}
	}
	log.Infof(ctx, "env: %s", app.Config().Environment)
	return app
}

func (a *App) initHTTP(_ context.Context) error {
	vars := a.Config().Variables

	http.NewServer(
		a.HTTPServer(),
		vars.GetString(appenv.SwaggerHost),
		vars.GetString(appenv.SwaggerScheme),
		// TODO
		nil,
	)
	return nil
}

func (a *App) initInfra(_ context.Context) error {
	vars := a.Config().Variables
	return nil
}

func (a *App) initRepositories(_ context.Context) error {
	return nil
}

// nolint
func (a *App) initServices(_ context.Context) error {
	vars := a.Config().Variables
	return nil
}

func (a *App) initClients(ctx context.Context) error {
	vars := a.Config().Variables

	a.smsClient = sms.NewClient(
		vars.GetString(appenv.SMSClientLogin),
		vars.GetString(appenv.SMSClientPassword),
	)

	a.whatsappClient = whatsapp.NewClient(
		vars.GetString(appenv.WhatsAppAPIURL),
		vars.GetString(appenv.WhatsAppMediaURL),
		vars.GetString(appenv.WhatsAppIDInstance),
		vars.GetString(appenv.WhatsAppAPIToken),
	)
	var err error
	a.s3Client, err = s3.NewClient(
		ctx,
		s3.Config{
			Region:          vars.GetString(appenv.AwsRegion),
			Endpoint:        vars.GetString(appenv.AwsS3Endpoint),
			AccessKeyID:     vars.GetString(appenv.AwsAccessKeyID),
			SecretAccessKey: vars.GetString(appenv.AwsSecretAccessKey),
			BucketName:      vars.GetString(appenv.AwsS3BucketName),
		},
	)

	return err
}

func (a *App) initJobs(_ context.Context) error {
	vars := a.Config().Variables

	// Загружаем таймзону
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}

	// Настройки cron с секундами и таймзоной
	a.WithCronOptions(
		cron.WithSeconds(),
		cron.WithLocation(loc),
	)

	a.WithCronJobs()

	return nil
}
