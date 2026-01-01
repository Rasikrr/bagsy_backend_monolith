package app

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/tokens"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/infra/jwt"
	formsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"
	networksR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/notifications"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/appenv"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	pointsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"
	usersR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	formsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	pointsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
)

type App struct {
	application.App

	smsClient      *sms.Client
	whatsappClient *whatsapp.Client

	tokensCache *tokens.Cache

	usersRepo    *usersR.Repository
	pointsRepo   *pointsR.Repository
	networksRepo *networksR.Repository
	formsRepo    *formsR.Repository

	usersService         *usersS.Service
	pointsService        *pointsS.Service
	networksService      *networksS.Service
	authService          *authS.Service
	formsService         *formsS.Service
	notificationsService *notifications.Service

	tokenManager *jwt.TokenManager
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initInfra,
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
	)
	return nil
}

func (a *App) initInfra(ctx context.Context) error {
	jwtSecret, err := a.Config().Variables.GetString(appenv.JWTSecret)
	if err != nil {
		return err
	}
	issuer, err := a.Config().Variables.GetString(appenv.JWTIssuer)
	if err != nil {
		return err
	}
	a.tokenManager = jwt.NewTokenManager(jwtSecret, issuer)
	return nil
}

func (a *App) initCache(_ context.Context) error {
	a.tokensCache = tokens.New(a.Redis())
	return nil
}

func (a *App) initRepositories(_ context.Context) error {
	a.usersRepo = usersR.NewRepository(a.Postgres())
	a.pointsRepo = pointsR.NewRepository(a.Postgres())
	a.formsRepo = formsR.NewRepository(a.Postgres())
	return nil
}

func (a *App) initServices(_ context.Context) error {
	accessTokenTTL, err := a.Config().Variables.GetDuration(appenv.AccessTokenTTL)
	if err != nil {
		return err
	}
	refreshTokenTTL, err := a.Config().Variables.GetDuration(appenv.RefreshTokenTTL)
	if err != nil {
		return err
	}
	registrationTokenTTL, err := a.Config().Variables.GetDuration(appenv.RegistrationTokenTTL)
	if err != nil {
		return err
	}
	registerConfirmationURL, err := a.Config().Variables.GetString(appenv.RegisterConfirmationURL)
	if err != nil {
		return err
	}

	a.usersService = usersS.NewService(a.usersRepo)
	a.networksService = networksS.NewService(a.networksRepo)
	a.pointsService = pointsS.NewService(a.pointsRepo, a.networksService)
	a.formsService = formsS.NewService(a.formsRepo)
	a.notificationsService = notifications.NewService(a.smsClient, a.whatsappClient, registerConfirmationURL)

	a.authService = authS.NewService(
		a.PostgresTXManager(),
		a.usersService,
		a.pointsService,
		a.notificationsService,
		a.tokenManager,
		a.tokensCache,
		a.tokensCache,
		accessTokenTTL,
		refreshTokenTTL,
		registrationTokenTTL,
	)

	return nil
}

func (a *App) initClients(_ context.Context) error {
	smsLogin, err := a.Config().Variables.GetString(appenv.SMSClientLogin)
	if err != nil {
		return err
	}
	smsPassword, err := a.Config().Variables.GetString(appenv.SMSClientPassword)
	if err != nil {
		return err
	}
	a.smsClient = sms.NewClient(smsLogin, smsPassword)

	whatsappAPIURL, err := a.Config().Variables.GetString(appenv.WhatsAppAPIURL)
	if err != nil {
		return err
	}
	whatsappAPIMediaURL, err := a.Config().Variables.GetString(appenv.WhatsAppMediaURL)
	if err != nil {
		return err
	}
	whatsappAPIIDInstance, err := a.Config().Variables.GetString(appenv.WhatsAppIDInstance)
	if err != nil {
		return err
	}
	whatsappAPIToken, err := a.Config().Variables.GetString(appenv.WhatsAppAPIToken)
	if err != nil {
		return err
	}

	a.whatsappClient = whatsapp.NewClient(
		whatsappAPIURL,
		whatsappAPIMediaURL,
		whatsappAPIIDInstance,
		whatsappAPIToken,
	)
	return nil
}

// nolint
func (a *App) initJobs(_ context.Context) error {
	////1 если хочешь добавить настройки (таймзону, кронджобы в секундах и т.д)
	//loc, err := time.LoadLocation("Asia/Almaty")
	//if err != nil {
	//	return err
	//}
	//
	//inactiveUserTTL, err := a.Config().Variables.GetDuration(appenv.InactiveUserTTL)
	//if err != nil {
	//	return err
	//}
	//inactiveUserJobSchedule, err := a.Config().Variables.GetString(appenv.InactiveUserJobSchedule)
	//if err != nil {
	//	return err
	//}
	//
	//a.WithCronOptions(cron.WithSeconds(), cron.WithLocation(loc))
	//// 2 доавбляешь джобы (расписание: если использовал опцию с секундами "* * * * * *", если не использовал, то "* * * * *")
	//// ВАЖНО: если используешь секунды, то все джобы должны быть в секундном формате
	//
	//a.WithCronJobs(
	//	jobs.NewExampleJob("example_job_2", "0 */1 * * * *"),
	//	jobs.NewDeleteUnactivatedUsers("delete_inactive_users", inactiveUserJobSchedule, inactiveUserTTL, a.usersService),
	//)

	return nil
}
