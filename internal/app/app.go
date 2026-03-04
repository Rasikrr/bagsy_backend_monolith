package app

import (
	"context"
	"time"

	appenv "github.com/Rasikrr/bagsy_backend_monolith/internal/appenvs"
	jwtinfra "github.com/Rasikrr/bagsy_backend_monolith/internal/infra/jwt"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/infra/messenger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	accessRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/access"
	actionTokenRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/action_token"
	otpRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/otp"
	bookingRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/booking"
	catalogRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/catalog"
	customerRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/customer"
	employeeRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/employee"
	locationRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/location"
	categoryRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/location_category"
	mediaRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media"
	orgRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/organization"
	planRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/plan"
	scheduleRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/schedule"
	subscriptionRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/subscription"
	workHistoryRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/work_history"

	pendingReg "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/pending_registraion"
	refreshTokenRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/tokens"
	invitePendingRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/invite/pending"

	authUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	bookingUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	inviteUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	locationUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	mediaUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/policy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/workers"

	"github.com/Rasikrr/bagsy_backend_monolith/pkg/s3"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/whatsapp"
	"github.com/Rasikrr/core/application"
	"github.com/Rasikrr/core/log"
	"github.com/robfig/cron/v3"
)

type App struct {
	application.App

	// Clients
	smsClient      *sms.Client
	whatsappClient *whatsapp.Client
	s3Client       *s3.Client

	// Repositories
	employeeRepo       *employeeRepo.Repository
	organizationRepo   *orgRepo.Repository
	planRepo           *planRepo.Repository
	subscriptionRepo   *subscriptionRepo.Repository
	workHistoryRepo    *workHistoryRepo.Repository
	accessRepo         *accessRepo.Repository
	categoryRepo       *categoryRepo.Repository
	locationRepo       *locationRepo.Repository
	bookingRepo        *bookingRepo.Repository
	customerRepo       *customerRepo.Repository
	catalogRepo        *catalogRepo.Repository
	scheduleRepo       *scheduleRepo.Repository
	mediaRepo          *mediaRepo.Repository
	otpRepo            *otpRepo.Repository
	pendingRegStore    *pendingReg.PendingRegistrationStore
	refreshTokenRepo   *refreshTokenRepo.RefreshTokenRepository
	actionTokenStore   *actionTokenRepo.Store
	pendingInviteStore *invitePendingRepo.Store

	// Infra
	tokenManager *jwtinfra.TokenManager
	tokenService *jwtinfra.TokenService
	messenger    *messenger.Messenger

	// Use Cases
	registerOwnerUC  *authUC.RegisterOwnerUseCase
	authUseCase      *authUC.UseCase
	resetPasswordUC  *authUC.ResetPasswordUseCase
	inviteEmployeeUC *inviteUC.UseCase
	employeeUseCase  *employeeUC.UseCase
	createLocationUC *locationUC.UseCase
	bookingUseCase   *bookingUC.UseCase
	mediaUseCase     *mediaUC.UseCase

	// Policies
	policy *policy.Policy
}

func InitApp(ctx context.Context) *App {
	app := &App{
		App: *application.NewApp(ctx),
	}
	for _, initFn := range []func(context.Context) error{
		app.initClients,
		app.initRepositories,
		app.initInfra,
		app.initUseCases,
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

func (a *App) initRepositories(_ context.Context) error {
	db := a.Postgres()

	a.employeeRepo = employeeRepo.NewRepository(db)
	a.organizationRepo = orgRepo.NewRepository(db)
	a.planRepo = planRepo.NewRepository(db)
	a.subscriptionRepo = subscriptionRepo.NewRepository(db)
	a.workHistoryRepo = workHistoryRepo.NewRepository(db)
	a.accessRepo = accessRepo.NewRepository(db)
	a.locationRepo = locationRepo.NewRepository(db)
	a.categoryRepo = categoryRepo.NewRepository(db)
	a.bookingRepo = bookingRepo.NewRepository(db)
	a.customerRepo = customerRepo.NewRepository(db)
	a.catalogRepo = catalogRepo.NewRepository(db)
	a.scheduleRepo = scheduleRepo.NewRepository(db)
	a.mediaRepo = mediaRepo.NewRepository(db)

	a.otpRepo = otpRepo.NewRepository(a.Redis())
	a.pendingRegStore = pendingReg.NewPendingRegistrationStore(a.Redis())
	a.refreshTokenRepo = refreshTokenRepo.NewRefreshTokenRepository(a.Redis())
	a.actionTokenStore = actionTokenRepo.NewStore(a.Redis())
	a.pendingInviteStore = invitePendingRepo.NewStore(a.Redis())

	return nil
}

func (a *App) initInfra(_ context.Context) error {
	vars := a.Config().Variables

	// JWT
	a.tokenManager = jwtinfra.NewTokenManager(
		vars.GetString(appenv.JWTSecret),
		vars.GetString(appenv.JWTIssuer),
	)

	accessTTL := vars.GetDuration(appenv.AccessTokenTTL)
	refreshTTL := vars.GetDuration(appenv.RefreshTokenTTL)

	a.tokenService = jwtinfra.NewTokenService(
		a.tokenManager,
		accessTTL,
		refreshTTL,
		a.refreshTokenRepo,
	)

	// Messenger (WhatsApp → SMS fallback)
	a.messenger = messenger.NewMessenger(a.whatsappClient, a.smsClient)

	return nil
}

func (a *App) initUseCases(_ context.Context) error {
	vars := a.Config().Variables
	txManager := a.PostgresTXManager()

	a.registerOwnerUC = authUC.NewRegisterOwnerUseCase(
		a.employeeRepo,
		a.planRepo,
		a.organizationRepo,
		a.subscriptionRepo,
		a.workHistoryRepo,
		a.tokenService,
		a.pendingRegStore,
		txManager,
		a.messenger,
	)

	a.authUseCase = authUC.NewUseCase(
		a.employeeRepo,
		a.employeeRepo,
		a.tokenService,
		a.actionTokenStore,
	)

	resetTTL := vars.GetDuration(appenv.PasswordResetTTL)
	frontendURL := vars.GetString(appenv.FrontendURL)

	a.resetPasswordUC = authUC.NewResetPasswordUseCase(
		a.employeeRepo,
		a.actionTokenStore,
		a.tokenService,
		a.messenger,
		resetTTL,
		frontendURL,
	)

	a.policy = policy.New()

	inviteTTL := vars.GetDuration(appenv.InviteTTL)
	inviteFrontendURL := vars.GetString(appenv.FrontendURL)

	a.inviteEmployeeUC = inviteUC.NewUseCase(
		a.employeeRepo,
		a.workHistoryRepo,
		a.actionTokenStore,
		a.pendingInviteStore,
		a.tokenService,
		a.messenger,
		a.policy,
		txManager,
		inviteTTL,
		inviteFrontendURL,
	)

	a.employeeUseCase = employeeUC.NewUseCase(
		a.employeeRepo,
		a.locationRepo,
		a.workHistoryRepo,
		a.mediaRepo,
		a.s3Client,
		txManager,
		a.policy,
		15*time.Minute,
	)

	a.createLocationUC = locationUC.NewCreateLocationUseCase(
		a.locationRepo,
		a.categoryRepo,
		a.organizationRepo,
		a.employeeRepo,
		a.policy,
		txManager,
	)

	mediaUploadTTL := vars.GetDuration(appenv.MediaUploadTTL)
	mediaMaxSize := int64(vars.GetInt(appenv.MediaMaxSizeBytes))

	a.mediaUseCase = mediaUC.NewUseCase(
		a.mediaRepo,
		a.s3Client,
		mediaUploadTTL,
		mediaMaxSize,
	)

	a.bookingUseCase = bookingUC.NewUseCase(
		a.bookingRepo,
		a.customerRepo,
		a.employeeRepo,
		a.catalogRepo,
		a.catalogRepo,
		a.locationRepo,
		a.subscriptionRepo,
		a.scheduleRepo,
		a.otpRepo,
		a.messenger,
		a.policy,
		txManager,
	)

	return nil
}

func (a *App) initHTTP(_ context.Context) error {
	vars := a.Config().Variables

	http.NewServer(
		a.HTTPServer(),
		vars.GetString(appenv.SwaggerHost),
		vars.GetString(appenv.SwaggerScheme),
		a.registerOwnerUC,
		a.authUseCase,
		a.resetPasswordUC,
		a.inviteEmployeeUC,
		a.employeeUseCase,
		a.accessRepo,
		a.createLocationUC,
		a.bookingUseCase,
		a.mediaUseCase,
	)
	return nil
}

func (a *App) initJobs(_ context.Context) error {
	vars := a.Config().Variables
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}

	a.WithCronOptions(
		cron.WithSeconds(),
		cron.WithLocation(loc),
	)

	mediaUploadTTL := vars.GetDuration(appenv.MediaUploadTTL)

	a.WithCronJobs(
		workers.NewMediaCleanupJob(a.mediaRepo, mediaUploadTTL, "0 */1 * * * *"),
	)

	return nil
}
