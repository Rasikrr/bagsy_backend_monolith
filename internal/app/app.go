package app

import (
	"context"

	bagsyconfirm "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/bagsy_confirm"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/register"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/tokens"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/s3"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/infra/jwt"
	formsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"
	mediaR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/media"
	pointMediaR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/point_media"
	userAvatarR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/user_avatar"
	networksR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	pointCategoriesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/point_categories"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/master_services"
	mediaS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/notifications"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/services"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/appenv"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	bagsiesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	masterServicesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/master_services"
	pointsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"
	servicesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/services"
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

	tokensCache       *tokens.Cache
	bagsyConfirmCache *bagsyconfirm.Cache
	registerCache     *register.Cache

	usersRepo           *usersR.Repository
	pointsRepo          *pointsR.Repository
	networksRepo        *networksR.Repository
	pointCategoriesRepo *pointCategoriesR.Repository
	formsRepo           *formsR.Repository
	bagsiesRepo         *bagsiesR.Repository
	masterServicesRepo  *masterServicesR.Repository
	servicesRepo        *servicesR.Repository
	mediaRepo           *mediaR.Repository
	userAvatarRepo      *userAvatarR.Repository
	pointMediaRepo      *pointMediaR.Repository

	usersService          *usersS.Service
	pointsService         *pointsS.Service
	networksService       *networksS.Service
	authService           *authS.Service
	formsService          *formsS.Service
	notificationsService  *notifications.Service
	bagsiesService        *bagsies.Service
	masterServicesService *masterservices.Service
	servicesService       *services.Service
	mediaService          *mediaS.Service

	s3Client *s3.Client

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
	vars := a.Config().Variables

	http.NewServer(
		a.HTTPServer(),
		a.Redis(),
		vars.GetString(appenv.SwaggerHost),
		vars.GetString(appenv.SwaggerScheme),
		a.authService,
		a.formsService,
		a.usersService,
		a.bagsiesService,
		a.pointsService,
		a.networksService,
		a.mediaService,
	)
	return nil
}

func (a *App) initInfra(_ context.Context) error {
	vars := a.Config().Variables

	a.tokenManager = jwt.NewTokenManager(
		vars.GetString(appenv.JWTSecret),
		vars.GetString(appenv.JWTIssuer),
	)
	return nil
}

func (a *App) initCache(_ context.Context) error {
	a.tokensCache = tokens.New(a.Redis())

	vars := a.Config().Variables

	a.bagsyConfirmCache = bagsyconfirm.NewCache(
		a.Redis(),
		vars.GetDuration(appenv.BagsyConfirmTTL),
	)

	a.registerCache = register.NewCache(
		a.Redis(),
	)
	return nil
}

func (a *App) initRepositories(_ context.Context) error {
	a.usersRepo = usersR.NewRepository(a.Postgres())
	a.pointsRepo = pointsR.NewRepository(a.Postgres())
	a.networksRepo = networksR.NewRepository(a.Postgres())
	a.pointCategoriesRepo = pointCategoriesR.NewRepository(a.Postgres())
	a.formsRepo = formsR.NewRepository(a.Postgres())
	a.masterServicesRepo = masterServicesR.NewRepository(a.Postgres())
	a.bagsiesRepo = bagsiesR.NewRepository(a.Postgres())
	a.servicesRepo = servicesR.NewRepository(a.Postgres())
	a.mediaRepo = mediaR.NewRepository(a.Postgres())
	a.userAvatarRepo = userAvatarR.NewRepository(a.Postgres())
	a.pointMediaRepo = pointMediaR.NewRepository(a.Postgres())
	return nil
}

func (a *App) initServices(_ context.Context) error {
	vars := a.Config().Variables

	a.networksService = networksS.NewService(a.networksRepo)

	a.mediaService = mediaS.NewService(
		a.PostgresTXManager(),
		a.s3Client,
		a.mediaRepo,
		a.userAvatarRepo,
		a.pointMediaRepo,
		vars.GetDuration(appenv.MediaTTL),
	)

	a.pointsService = pointsS.NewService(
		a.pointsRepo,
		a.networksService,
		a.pointCategoriesRepo,
		a.mediaService,
		a.PostgresTXManager(),
	)
	a.formsService = formsS.NewService(a.formsRepo)
	a.notificationsService = notifications.NewService(
		a.smsClient,
		a.whatsappClient,
		vars.GetString(appenv.RegisterConfirmationURL),
	)
	a.masterServicesService = masterservices.NewService(a.masterServicesRepo)
	a.servicesService = services.NewService(a.servicesRepo)

	a.usersService = usersS.NewService(
		a.PostgresTXManager(),
		a.usersRepo,
		a.pointsService,
		a.mediaService,
	)

	a.authService = authS.NewService(
		a.PostgresTXManager(),
		a.usersService,
		a.pointsService,
		a.notificationsService,
		a.networksService,
		a.tokenManager,
		a.tokensCache,
		a.tokensCache,
		a.registerCache,
		vars.GetDuration(appenv.AccessTokenTTL),
		vars.GetDuration(appenv.RefreshTokenTTL),
		vars.GetDuration(appenv.RegistrationTTL),
	)

	a.bagsiesService = bagsies.NewService(
		a.PostgresTXManager(),
		a.bagsiesRepo,
		a.masterServicesService,
		a.servicesService,
		a.usersService,
		a.notificationsService,
		a.bagsyConfirmCache,
	)

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
