package app

import (
	"context"
	"time"

	networksR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/telegram"
	"github.com/robfig/cron/v3"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/appenv"
	authC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/auth"
	smsC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/jobs"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	bagsiesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	formsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"
	pointsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"
	usersR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	bagsiesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	formsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	pointsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
)

type App struct {
	application.App

	smsCache          smsC.Cache
	authCache         authC.Cache
	bagsiesCodeCacche authC.Cache

	smsClient      sms.Client
	tgClient       telegram.Client
	whatsAppClient whatsapp.Client

	usersRepo    usersR.Repository
	bagsiesRepo  bagsiesR.Repository
	formsRepo    formsR.Repository
	pointsRepo   pointsR.Repository
	networksRepo networksR.Repository

	authService     authS.Service
	formsService    formsS.Service
	usersService    usersS.Service
	pointsService   pointsS.Service
	bagsiesService  bagsiesS.Service
	networksService networksS.Service
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
		a.formsService,
		a.usersService,
		a.bagsiesService,
		a.pointsService,
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

	bagsiesCodeTTL, err := a.Config().Variables.GetDuration(appenv.BagsiesCodeTTL)
	if err != nil {
		return err
	}

	a.authCache = authC.NewCache(a.Redis(), authCodeTTL)
	a.bagsiesCodeCacche = authC.NewCache(a.Redis(), bagsiesCodeTTL)
	return nil
}

func (a *App) initRepositories(_ context.Context) error {
	a.usersRepo = usersR.NewRepository(a.Postgres())
	a.bagsiesRepo = bagsiesR.NewRepository(a.Postgres())
	a.formsRepo = formsR.NewRepository(a.Postgres())
	a.pointsRepo = pointsR.NewRepository(a.Postgres())
	a.networksRepo = networksR.NewRepository(a.Postgres())
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
	accessTokenTTL, err := a.Config().Variables.GetDuration(appenv.AccessTokenTTL)
	if err != nil {
		return err
	}
	refreshTokenTTL, err := a.Config().Variables.GetDuration(appenv.RefreshTokenTTL)
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
		accessTokenTTL,
		refreshTokenTTL,
	)

	a.bagsiesService = bagsiesS.NewService(
		a.whatsAppClient,
		a.bagsiesCodeCacche,
		a.bagsiesRepo,
		a.usersRepo,
		a.PostgresTXManager(),
	)

	a.formsService = formsS.NewService(a.formsRepo)

	a.pointsService = pointsS.NewService(a.pointsRepo, a.networksRepo)

	a.networksService = networksS.NewService(a.networksRepo)

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

// nolint
func (a *App) initJobs(_ context.Context) error {
	//1 если хочешь добавить настройки (таймзону, кронджобы в секундах и т.д)
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}

	inactiveUserTTL, err := a.Config().Variables.GetDuration(appenv.InactiveUserTTL)
	if err != nil {
		return err
	}
	inactiveUserJobSchedule, err := a.Config().Variables.GetString(appenv.InactiveUserJobSchedule)
	if err != nil {
		return err
	}

	a.WithCronOptions(cron.WithSeconds(), cron.WithLocation(loc))
	// 2 доавбляешь джобы (расписание: если использовал опцию с секундами "* * * * * *", если не использовал, то "* * * * *")
	// ВАЖНО: если используешь секунды, то все джобы должны быть в секундном формате

	a.WithCronJobs(
		jobs.NewExampleJob("example_job_2", "0 */1 * * * *"),
		jobs.NewDeleteUnactivatedUsers("delete_inactive_users", inactiveUserJobSchedule, inactiveUserTTL, a.usersService),
	)

	return nil
}
