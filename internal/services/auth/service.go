package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
)

const (
	maxVerificationAttempts = 3
)

type usersService interface {
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	UpdatePasswordByPhone(ctx context.Context, phone, rawPassword string) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*point.Point, error)
}

type registrationService interface {
	RegisterNewOwner(ctx context.Context, cmd *auth.RegisterManagementCommand) (*user.User, error)
	RegisterNewStaff(ctx context.Context, cmd *auth.RegisterStaffCommand, rawPassword string) (*user.User, error)
}

type notificationService interface {
	SendStaffRegistrationLink(ctx context.Context, phone, token string) error
	SendManagementAuthConfirmationCode(ctx context.Context, phone, code string) error
	SendPasswordChangeLink(ctx context.Context, phone, token string) error
}

type tokenManager interface {
	NewAccessToken(payload *AccessTokenPayload, ttl time.Duration) (string, error)
	NewRefreshToken() (raw, hash string, err error)
	ParseAccessToken(token string) (*AccessTokenPayload, error)
}

type registerCache interface {
	SaveManagementRequest(ctx context.Context, state *ManagementRegistrationState, duration time.Duration) error
	GetManagementRequest(ctx context.Context, phone string) (*ManagementRegistrationState, error)
	DeleteManagementRequest(ctx context.Context, phone string) error

	SaveStaffRequest(ctx context.Context, req *auth.RegisterStaffCommand, duration time.Duration) error
	GetStaffRequest(ctx context.Context, phone string) (*auth.RegisterStaffCommand, error)
	DeleteStaffRequest(ctx context.Context, phone string) error
}

// refreshTokenRepository управляет хранением refresh токенов
// Хранит маппинг tokenHash → phone для валидации токенов
// Это позволяет поддерживать multiple devices (несколько токенов на пользователя)
type refreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, tokenHash, phone string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, tokenHash string) (string, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}

// registrationTokenCache управляет one-time use registration токенами
type tokensCache interface {
	SaveInviteToken(ctx context.Context, token string, payload *InviteTokenInfo, ttl time.Duration) error
	GetInviteToken(ctx context.Context, token string) (*InviteTokenInfo, error)
	DeleteInviteToken(ctx context.Context, token string) error
}

type Service struct {
	txManager              database.TXManager
	registrationService    registrationService
	usersService           usersService
	pointsService          pointsService
	notificationService    notificationService
	tokenManager           tokenManager
	refreshTokenRepository refreshTokenRepository
	tokensCache            tokensCache
	registerCache          registerCache
	accessTTL              time.Duration
	refreshTTL             time.Duration
	registrationTTL        time.Duration
}

func NewService(
	txManager database.TXManager,
	registrationService registrationService,
	usersService usersService,
	pointsService pointsService,
	notificationService notificationService,
	tokenManager tokenManager,
	refreshTokenRepo refreshTokenRepository,
	tokensCache tokensCache,
	registerCache registerCache,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	registrationTTL time.Duration,
) *Service {
	return &Service{
		txManager:              txManager,
		registrationService:    registrationService,
		usersService:           usersService,
		pointsService:          pointsService,
		notificationService:    notificationService,
		tokenManager:           tokenManager,
		refreshTokenRepository: refreshTokenRepo,
		tokensCache:            tokensCache,
		registerCache:          registerCache,
		accessTTL:              accessTTL,
		refreshTTL:             refreshTTL,
		registrationTTL:        registrationTTL,
	}
}

func (s *Service) Login(ctx context.Context, phone string, password string) (access, refresh string, err error) {
	user, err := s.usersService.GetByPhone(ctx, phone)
	if err != nil {
		if domainErr.IsNotFound(err) {
			return "", "", auth.ErrInvalidCredentials
		}
		return "", "", err
	}

	if !user.Active {
		return "", "", auth.ErrUserInactive
	}

	valid := hash.CheckPassword(user.PasswordHash, password)
	if !valid {
		return "", "", auth.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", err // Уже обернута в domain error
	}

	return accessToken, refreshToken, nil
}

func (s *Service) InspectAuthToken(ctx context.Context, token string) (*InviteTokenInfo, error) {
	payload, err := s.tokensCache.GetInviteToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *Service) RegisterManagement(ctx context.Context, cmd *auth.RegisterManagementCommand) error {
	// Проверка, что запрос не был уже создан ранее
	existedReq, err := s.registerCache.GetManagementRequest(ctx, cmd.Phone)
	if err != nil {
		if !domainErr.IsNotFound(err) {
			return err
		}
	}
	if existedReq != nil {
		return domainErr.NewConflictError("registration request already created", nil)
	}

	authCode := codegen.GenerateAuthCode()
	registerState := newManagementRegistrationState(cmd, authCode)

	err = s.registerCache.SaveManagementRequest(ctx, registerState, s.registrationTTL)
	if err != nil {
		return err
	}
	return s.notificationService.SendManagementAuthConfirmationCode(ctx, cmd.Phone, authCode)
}

func (s *Service) ResendRegisterManagementCode(ctx context.Context, phone string) error {
	state, err := s.registerCache.GetManagementRequest(ctx, phone)
	if err != nil {
		return err
	}
	newCode := codegen.GenerateAuthCode()
	state.AuthCode = newCode
	state.Attempts = 0
	err = s.registerCache.SaveManagementRequest(ctx, state, s.registrationTTL)
	if err != nil {
		return err
	}
	return s.notificationService.SendManagementAuthConfirmationCode(ctx, phone, newCode)
}

func (s *Service) RegisterManagementConfirm(ctx context.Context, phone, code string) (access, refresh string, err error) {
	state, err := s.registerCache.GetManagementRequest(ctx, phone)
	if err != nil {
		return "", "", err
	}
	if state.AuthCode != code {
		state.Attempts++
		if state.Attempts >= maxVerificationAttempts {
			// Удаляем данные регистрации после превышения лимита попыток
			if delErr := s.registerCache.DeleteManagementRequest(ctx, phone); delErr != nil {
				return "", "", delErr
			}
			return "", "", auth.ErrTooManyVerificationAttempts
		}
		// Сохраняем обновленный счетчик попыток
		if saveErr := s.registerCache.SaveManagementRequest(ctx, state, s.registrationTTL); saveErr != nil {
			return "", "", saveErr
		}
		return "", "", auth.ErrInvalidVerificationCode.
			WithDetail("attempts_remaining", maxVerificationAttempts-state.Attempts)
	}
	var newUser *user.User

	newUser, err = s.registrationService.RegisterNewOwner(ctx, state.Command)

	if err != nil {
		return "", "", err
	}
	return s.generateTokens(ctx, newUser)
}

// RegisterStaff регистрирует нового работника (orchestrator)
// Проверяет права, создает user, генерирует и отправляет ссылку для подтверждения
func (s *Service) RegisterStaff(ctx context.Context, cmd *auth.RegisterStaffCommand) error {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return err
	}

	// 1. Валидация прав доступа
	if validErr := s.validateStaffRegistrationPermissions(ctx, cmd.PointCode, cmd.Role, act); validErr != nil {
		return validErr
	}
	cmd.NetworkCode = act.NetworkCode()

	// 2. Проверим что запрос не был создан ранее
	existedReq, err := s.registerCache.GetStaffRequest(ctx, cmd.Phone)
	if err != nil {
		if !domainErr.IsNotFound(err) {
			return err
		}
	}
	if existedReq != nil {
		return domainErr.NewConflictError("registration request already created", nil)
	}

	// 3. Сохраним данные в кеше чтобы преждевременно не создавать юзера
	err = s.registerCache.SaveStaffRequest(ctx, cmd, s.registrationTTL)
	if err != nil {
		return err
	}

	// 4. Генерируем короткий invite токен (8-10 символов)
	inviteToken := codegen.GenerateAuthToken()

	payload := &InviteTokenInfo{
		Phone:       cmd.Phone,
		PointCode:   cmd.PointCode,
		NetworkCode: cmd.NetworkCode,
		Purpose:     TokenPurposeRegister,
	}

	err = s.tokensCache.SaveInviteToken(ctx, inviteToken, payload, s.registrationTTL)
	if err != nil {
		return err
	}

	// 5. Отправляем уведомление (WhatsApp с fallback на SMS)
	if sendErr := s.notificationService.SendStaffRegistrationLink(ctx, cmd.Phone, inviteToken); sendErr != nil {
		return sendErr
	}
	return nil
}

func (s *Service) ResendRegisterStaffLink(ctx context.Context, phone string) error {
	cmd, err := s.registerCache.GetStaffRequest(ctx, phone)
	if err != nil {
		return err
	}

	// Генерируем новый короткий invite токен
	token := codegen.GenerateAuthToken()
	payload := &InviteTokenInfo{
		Phone:       cmd.Phone,
		PointCode:   cmd.PointCode,
		NetworkCode: cmd.NetworkCode,
	}

	if saveErr := s.tokensCache.SaveInviteToken(ctx, token, payload, s.registrationTTL); saveErr != nil {
		return saveErr
	}

	// ВАЖНО: Пересохраняем данные в кэш, чтобы обновить TTL!
	// Теперь данные будут жить еще registrationTTL (24 часа) от момента resend,
	// синхронизируя время жизни кэша с временем жизни токена
	if saveErr := s.registerCache.SaveStaffRequest(ctx, cmd, s.registrationTTL); saveErr != nil {
		return domainErr.NewInternalError("failed to refresh staff request TTL", saveErr)
	}

	return s.notificationService.SendStaffRegistrationLink(ctx, cmd.Phone, token)
}

// nolint: govet
func (s *Service) RegisterStaffConfirm(ctx context.Context, confirmCmd *auth.RegisterStaffConfirmCommand) (access, refresh string, err error) {
	// 1. Получаем данные из Redis по короткому токену
	tokenPayload, err := s.tokensCache.GetInviteToken(ctx, confirmCmd.Token)
	if err != nil {
		return "", "", domainErr.NewUnauthorizedError("invalid or expired registration token").WithError(err)
	}

	cmd, err := s.registerCache.GetStaffRequest(ctx, tokenPayload.Phone)
	if err != nil {
		return "", "", err
	}

	// 2. Регистрируем и привязываем юзера
	newStaff, err := s.registrationService.RegisterNewStaff(ctx, cmd, confirmCmd.Password)
	if err != nil {
		return "", "", err
	}
	// 3. Удаляем использованный invite токен (игнорируем ошибки - токен истечет по TTL)
	_ = s.tokensCache.DeleteInviteToken(ctx, confirmCmd.Token)

	return s.generateTokens(ctx, newStaff)
}

func (s *Service) SendPasswordChangeLink(ctx context.Context, phone string) error {
	exist, err := s.usersService.ExistsByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if !exist {
		return user.ErrUserNotFound
	}
	// Генерируем короткий auth токен
	token := codegen.GenerateAuthToken()
	payload := &InviteTokenInfo{
		Phone:   phone,
		Purpose: TokenPurposePasswordChange,
	}

	if saveErr := s.tokensCache.SaveInviteToken(ctx, token, payload, s.registrationTTL); saveErr != nil {
		return saveErr
	}

	return s.notificationService.SendPasswordChangeLink(ctx, phone, token)
}

func (s *Service) ChangePassword(ctx context.Context, req *auth.ChangePasswordConfirmCommand) error {
	// 1. Получаем данные из Redis по короткому токену
	tokenPayload, err := s.tokensCache.GetInviteToken(ctx, req.Token)
	if err != nil {
		return domainErr.NewUnauthorizedError("invalid or expired token").WithError(err)
	}

	err = s.usersService.UpdatePasswordByPhone(ctx, tokenPayload.Phone, req.Password)
	if err != nil {
		return err
	}

	// 2. Удаляем использованный токен (игнорируем ошибки - токен истечет по TTL)
	_ = s.tokensCache.DeleteInviteToken(ctx, req.Token)

	return nil
}

// validateStaffRegistrationPermissions проверяет права на регистрацию работника
func (s *Service) validateStaffRegistrationPermissions(ctx context.Context, pointCode string, role user.Role, createdBy *actor.Actor) error {
	switch createdBy.Role() {
	case user.RoleStaff:
		// Staff не может создавать пользователей
		return domainErr.NewForbiddenError("staff cannot create users")

	case user.RoleManager:
		// Point Manager может создавать только в своей точке
		if !role.OneOf(user.RoleStaff) {
			return domainErr.NewForbiddenError("manager cannot create (net)managers")
		}
		if pointCode != createdBy.PointCode() {
			return domainErr.NewForbiddenError("point manager can only create staff in their own point").
				WithDetail("allowed_point", createdBy.PointCode()).
				WithDetail("requested_point", pointCode)
		}

	case user.RoleNetManager, user.RoleSelfOwner:
		if !role.OneOf(user.RoleManager, user.RoleStaff) {
			return domainErr.NewForbiddenError("manager cannot create net managers")
		}
		// Network Manager и SelfOwner могут создавать в любой точке своей сети
		// Дополнительно проверяем, что точка принадлежит сети
		point, err := s.pointsService.GetByCode(ctx, pointCode)
		if err != nil {
			return err
		}

		if point.NetworkCode != createdBy.NetworkCode() {
			return domainErr.NewForbiddenError("point does not belong to the network").
				WithDetail("point_network", point.NetworkCode).
				WithDetail("requested_network", createdBy.NetworkCode())
		}

	default:
		return domainErr.NewForbiddenError("unknown role").
			WithDetail("role", createdBy.Role().String())
	}

	return nil
}

// Logout удаляет refresh токен пользователя
// Access токен не blacklist'им, т.к. он короткоживущий (15 мин)
// Если нужна немедленная инвалидация access токена, можно добавить blacklist позже
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// Хэшируем refresh токен
	h := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(h[:])

	// Удаляем refresh токен из cache
	// Игнорируем ошибку если токен уже не существует
	if err := s.refreshTokenRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return domainErr.NewInternalError("failed to delete refresh token", err)
	}

	return nil
}

// VerifyAccessToken парсит access токен и создает domain.Actor
// Это метод согласно Clean Architecture (JWT → DTO → Actor)
// Service не зависит от конкретных ошибок JWT - просто оборачивает любую ошибку
func (s *Service) VerifyAccessToken(_ context.Context, tokenStr string) (*actor.Actor, error) {
	// 1. Парсим токен через JWT (получаем DTO из domain)
	payload, err := s.tokenManager.ParseAccessToken(tokenStr)
	if err != nil {
		// 2. Любая ошибка парсинга → Unauthorized
		// Детали ошибки сохраняются через wrapping для логирования
		return nil, domainErr.NewUnauthorizedError("authentication failed").WithError(err)
	}

	// 3. Валидируем роль (domain logic!)
	role, err := user.RoleString(payload.Role)
	if err != nil {
		return nil, auth.ErrInvalidTokenClaims.
			WithDetail("claim", "role").
			WithError(err)
	}

	// 4. Создаем Session (domain entity)
	act := actor.NewActor().
		SetPhone(payload.Phone).
		SetRole(role).
		SetPointCode(payload.PointCode).
		SetNetworkCode(payload.NetworkCode)

	return act, nil
}

// Refresh обновляет пару токенов по refresh токену
// Использует hash-based подход:
// 1. Хэширует raw токен от клиента
// 2. Ищет hash в cache → получает phone
// 3. Удаляет старый токен (ротация)
// 4. Генерирует новую пару (автоматически сохраняет новый hash)
func (s *Service) Refresh(ctx context.Context, rawToken string) (string, string, error) {
	// 1. Хэшируем refresh токен (тот же алгоритм, что в tokenManager.NewRefreshToken)
	h := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(h[:])

	// 2. Проверяем наличие токена в cache и получаем phone
	phone, err := s.refreshTokenRepository.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		// GetRefreshToken возвращает domain error если токен не найден/истек
		return "", "", err
	}

	// 3. Удаляем старый refresh токен (ротация для безопасности)
	// Делаем это ДО получения пользователя, чтобы токен нельзя было использовать повторно
	if delErr := s.refreshTokenRepository.DeleteRefreshToken(ctx, tokenHash); delErr != nil {
		return "", "", domainErr.NewInternalError("failed to delete old refresh token", delErr)
	}

	// 4. Получаем актуальные данные пользователя из БД
	user, err := s.usersService.GetByPhone(ctx, phone)
	if err != nil {
		if domainErr.IsNotFound(err) {
			return "", "", auth.ErrInvalidCredentials
		}
		return "", "", domainErr.NewInternalError("failed to get user", err)
	}

	// 5. Проверяем активность пользователя
	if !user.Active {
		return "", "", auth.ErrUserInactive
	}

	// 6. Генерируем новую пару токенов (автоматически сохраняет новый refresh hash)
	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) generateTokens(ctx context.Context, user *user.User) (accessToken, refreshToken string, err error) {
	payload := newAccessTokenPayload(
		user.Phone,
		user.Role.String(),
		ptr.Deref(user.PointCode),
		ptr.Deref(user.NetworkCode),
	)
	access, err := s.tokenManager.NewAccessToken(payload, s.accessTTL)
	if err != nil {
		return "", "", domainErr.NewInternalError("failed to generate access token", err)
	}
	refresh, refreshHash, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", "", domainErr.NewInternalError("failed to generate refresh token", err)
	}

	// Сохраняем hash refresh токена в cache (tokenHash → phone)
	if saveErr := s.refreshTokenRepository.SaveRefreshToken(ctx, refreshHash, user.Phone, s.refreshTTL); saveErr != nil {
		return "", "", domainErr.NewInternalError("failed to save refresh token", saveErr)
	}

	return access, refresh, nil
}
