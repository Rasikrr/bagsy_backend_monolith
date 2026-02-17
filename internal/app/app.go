package app

import (
	"context"
	"time"

	appenv "github.com/Rasikrr/bagsy_backend_monolith/internal/appenvs"
	jwtinfra "github.com/Rasikrr/bagsy_backend_monolith/internal/infra/jwt"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/infra/messenger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http"
	employeeRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/employee"
	orgRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/organization"
	planRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/plan"
	subscriptionRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/subscription"
	workHistoryRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/work_history"

	pendingReg "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/pending_registraion"
	refreshTokenRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/auth/tokens"

	authUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"

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
	employeeRepo     *employeeRepo.Repository
	organizationRepo *orgRepo.Repository
	planRepo         *planRepo.Repository
	subscriptionRepo *subscriptionRepo.Repository
	workHistoryRepo  *workHistoryRepo.Repository
	pendingRegStore  *pendingReg.PendingRegistrationStore
	refreshTokenRepo *refreshTokenRepo.RefreshTokenRepository

	// Infra
	tokenManager *jwtinfra.TokenManager
	tokenService *jwtinfra.TokenService
	otpSender    *messenger.OTPSender

	// Use Cases
	registerOwnerUC *authUC.RegisterOwnerUseCase
	authUseCase     *authUC.UseCase
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

	a.pendingRegStore = pendingReg.NewPendingRegistrationStore(a.Redis())
	a.refreshTokenRepo = refreshTokenRepo.NewRefreshTokenRepository(a.Redis())

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

	// OTP Sender (WhatsApp → SMS fallback)
	a.otpSender = messenger.NewOTPSender(a.whatsappClient, a.smsClient)

	return nil
}

func (a *App) initUseCases(_ context.Context) error {
	txm := a.PostgresTXManager()

	a.registerOwnerUC = authUC.NewRegisterOwnerUseCase(
		a.employeeRepo,
		a.planRepo,
		a.organizationRepo,
		a.subscriptionRepo,
		a.workHistoryRepo,
		a.tokenService,
		a.pendingRegStore,
		txm,
		a.otpSender,
	)

	a.authUseCase = authUC.NewUseCase(a.employeeRepo, a.tokenService)

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
	)
	return nil
}

func (a *App) initJobs(_ context.Context) error {
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}

	a.WithCronOptions(
		cron.WithSeconds(),
		cron.WithLocation(loc),
	)

	a.WithCronJobs()

	return nil
}
