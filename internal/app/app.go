package app

import (
	"context"
	"time"

	bagsyconfirm "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/bagsy_confirm"
	pointCategoriesC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/point_categories"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/register"
	serviceCategoriesC "github.com/Rasikrr/bagsy_backend_monolith/internal/cache/service_categories"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/tokens"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/s3"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/infra/jwt"
	bagsyNotificationsJob "github.com/Rasikrr/bagsy_backend_monolith/internal/jobs/bagsy_notifications"
	formsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"
	mediaR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/media"
	pointMediaR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/point_media"
	userAvatarR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/user_avatar"
	networksR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	notificationsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/notifications"
	pointCategoriesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/point_categories"
	pointCategoryServicesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/point_category_services"
	serviceCategoriesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/service_categories"
	serviceSubcategoryR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/service_subcategory"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	bagsyNotificationsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsy_notifications"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/master_services"
	mediaS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/notifications"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/registration"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/services"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
	"github.com/robfig/cron/v3"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/appenv"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	bagsiesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	masterServicesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/master_services"
	pointsR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"
	servicesR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/services"
	usersR "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	formsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/media/point_photos"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/media/users_photos"

	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	pointCategoriesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/point_categories"
	pointsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	serviceCategoriesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/service_categories"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
)

type App struct {
	application.App

	smsClient      *sms.Client
	whatsappClient *whatsapp.Client

	tokensCache            *tokens.Cache
	bagsyConfirmCache      *bagsyconfirm.Cache
	registerCache          *register.Cache
	pointCategoriesCache   *pointCategoriesC.Cache
	serviceCategoriesCache *serviceCategoriesC.Cache

	usersRepo                 *usersR.Repository
	pointsRepo                *pointsR.Repository
	networksRepo              *networksR.Repository
	pointCategoriesRepo       *pointCategoriesR.Repository
	formsRepo                 *formsR.Repository
	bagsiesRepo               *bagsiesR.Repository
	masterServicesRepo        *masterServicesR.Repository
	servicesRepo              *servicesR.Repository
	mediaRepo                 *mediaR.Repository
	userAvatarRepo            *userAvatarR.Repository
	pointMediaRepo            *pointMediaR.Repository
	notificationsRepo         *notificationsR.Repository
	pointCategoryServicesRepo *pointCategoryServicesR.Repository
	serviceCategoriesRepo     *serviceCategoriesR.Repository
	serviceSubcategoriesRepo  *serviceSubcategoryR.Repository

	usersService              *usersS.Service
	pointsService             *pointsS.Service
	networksService           *networksS.Service
	authService               *authS.Service
	formsService              *formsS.Service
	notificationsService      *notifications.Service
	bagsyNotificationsService *bagsyNotificationsS.Service
	bagsiesService            *bagsies.Service
	masterServicesService     *masterservices.Service
	servicesService           *services.Service
	mediaService              *mediaS.Service
	pointsMediaService        *pointphotos.Service
	userPhotosService         *usersphotos.Service
	registrationService       *registration.Service
	pointCategoriesService    *pointCategoriesS.Service
	serviceCategoriesService  *serviceCategoriesS.Service

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
		a.servicesService,
		a.mediaService,
		a.pointCategoriesService,
		a.serviceCategoriesService,
		a.masterServicesService,
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

	a.bagsyConfirmCache = bagsyconfirm.NewCache(
		a.Redis(),
	)

	a.registerCache = register.NewCache(
		a.Redis(),
	)

	a.pointCategoriesCache = pointCategoriesC.New(a.Redis())
	a.serviceCategoriesCache = serviceCategoriesC.New(a.Redis())
	return nil
}

func (a *App) initRepositories(_ context.Context) error {
	a.usersRepo = usersR.NewRepository(a.Postgres())
	a.pointsRepo = pointsR.NewRepository(a.Postgres())
	a.networksRepo = networksR.NewRepository(a.Postgres())
	a.pointCategoriesRepo = pointCategoriesR.NewRepository(a.Postgres())
	a.pointCategoryServicesRepo = pointCategoryServicesR.NewRepository(a.Postgres())
	a.serviceCategoriesRepo = serviceCategoriesR.NewRepository(a.Postgres())
	a.serviceSubcategoriesRepo = serviceSubcategoryR.NewRepository(a.Postgres())
	a.formsRepo = formsR.NewRepository(a.Postgres())
	a.masterServicesRepo = masterServicesR.NewRepository(a.Postgres())
	a.bagsiesRepo = bagsiesR.NewRepository(a.Postgres())
	a.servicesRepo = servicesR.NewRepository(a.Postgres())
	a.mediaRepo = mediaR.NewRepository(a.Postgres())
	a.userAvatarRepo = userAvatarR.NewRepository(a.Postgres())
	a.pointMediaRepo = pointMediaR.NewRepository(a.Postgres())
	a.notificationsRepo = notificationsR.NewRepository(a.Postgres())
	return nil
}

// nolint
func (a *App) initServices(_ context.Context) error {
	vars := a.Config().Variables

	a.networksService = networksS.NewService(a.networksRepo)

	a.pointCategoriesService = pointCategoriesS.NewService(
		a.pointCategoriesRepo,
		a.pointCategoriesCache,
		vars.GetDuration(appenv.PointCategoriesTTL),
	)

	a.mediaService = mediaS.NewService(
		a.PostgresTXManager(),
		a.s3Client,
		a.mediaRepo,
		vars.GetDuration(appenv.MediaTTL),
	)
	a.pointsMediaService = pointphotos.NewService(
		a.PostgresTXManager(),
		a.pointMediaRepo,
		a.mediaService,
		vars.GetInt(appenv.PointMediaMaxCount),
	)
	a.userPhotosService = usersphotos.NewService(
		a.PostgresTXManager(),
		a.userAvatarRepo,
		a.mediaService,
	)

	a.serviceCategoriesService = serviceCategoriesS.NewService(
		a.pointsService,
		a.pointCategoryServicesRepo,
		a.serviceCategoriesRepo,
		a.serviceSubcategoriesRepo,
		a.serviceCategoriesCache,
		vars.GetDuration(appenv.ServiceCategoriesTTL),
	)

	a.formsService = formsS.NewService(a.formsRepo)
	a.notificationsService = notifications.NewService(
		a.smsClient,
		a.whatsappClient,
		vars.GetString(appenv.RegisterConfirmationURL),
	)
	a.masterServicesService = masterservices.NewService(a.masterServicesRepo, a.usersRepo, a.servicesRepo)
	a.servicesService = services.NewService(
		a.servicesRepo,
		a.masterServicesRepo,
		a.serviceCategoriesRepo,
		a.serviceSubcategoriesRepo,
	)

	// Сервис планирования уведомлений о записях
	a.bagsyNotificationsService = bagsyNotificationsS.NewService(
		a.notificationsRepo,
		vars.GetInt(appenv.BagsyNotificationMaxAttempts),
	)

	a.usersService = usersS.NewService(
		a.PostgresTXManager(),
		a.usersRepo,
		a.pointsService,
		a.userPhotosService,
	)

	a.pointsService = pointsS.NewService(
		a.pointsRepo,
		a.networksService,
		a.pointCategoriesRepo,
		a.pointsMediaService,
		a.usersService,
		a.PostgresTXManager(),
	)

	a.registrationService = registration.NewService(
		a.PostgresTXManager(),
		a.usersService,
		a.networksService,
	)

	a.authService = authS.NewService(
		a.PostgresTXManager(),
		a.registrationService,
		a.usersService,
		a.pointsService,
		a.notificationsService,
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
		a.pointsService,
		a.notificationsService,
		a.bagsyNotificationsService,
		a.bagsyConfirmCache,
		vars.GetDuration(appenv.BagsyConfirmTTL),
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

func (a *App) initJobs(_ context.Context) error {
	vars := a.Config().Variables

	// Загружаем таймзону
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}

	// Настройки cron с секундами и таймзоной
	a.WithCronOptions(cron.WithSeconds(), cron.WithLocation(loc))

	// Создаем адаптер для отправки уведомлений
	messengerAdapter := bagsyNotificationsJob.NewMessengerAdapter(
		a.notificationsService,
		a.servicesService,
		a.usersService,
		a.pointsService,
		loc,
	)

	// Создаем джобу для отправки уведомлений о записях
	notificationsJob := bagsyNotificationsJob.NewJob(
		"bagsy_notifications",
		vars.GetString(appenv.BagsyNotificationSchedule), // например "0 */1 * * * *" (каждую минуту)
		vars.GetInt(appenv.BagsyNotificationBatchSize),   // например 100
		vars.GetInt(appenv.BagsyNotificationWorkerCount), // например 5
		a.bagsyNotificationsService,
		a.bagsiesService,
		messengerAdapter,
	)

	a.WithCronJobs(notificationsJob)

	return nil
}
