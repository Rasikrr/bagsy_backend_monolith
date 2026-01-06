package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
)

type usersService interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
}

type notificationService interface {
	SendRegistrationLink(ctx context.Context, phone, token string) error
}

type tokenManager interface {
	NewAccessToken(user *entity.User, ttl time.Duration) (string, error)
	NewRefreshToken() (raw, hash string, err error)
	NewRegistrationToken(dto *dto.RegistrationTokenPayload, ttl time.Duration) (string, error)
	ParseAccessToken(token string) (*dto.AccessTokenPayload, error)
	ParseRegistrationToken(token string) (*dto.RegistrationTokenPayload, error)
}

// refreshTokenRepository управляет хранением refresh токенов
// Хранит маппинг tokenHash → phone для валидации токенов
// Это позволяет поддерживать multiple devices (несколько токенов на пользователя)
type refreshTokenRepository interface {
	// SaveRefreshToken сохраняет hash refresh токена
	// Key: tokenHash, Value: phone
	SaveRefreshToken(ctx context.Context, tokenHash, phone string, ttl time.Duration) error

	// GetRefreshToken получает phone по hash токена
	// Возвращает domain error если токен не найден
	GetRefreshToken(ctx context.Context, tokenHash string) (string, error)

	// DeleteRefreshToken удаляет refresh токен по hash
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}

// registrationTokenCache управляет one-time use registration токенами
type registrationTokenCache interface {
	// MarkRegistrationTokenAsUsed помечает токен как использованный (атомарно через SET NX)
	// Возвращает true если токен уже был использован ранее
	MarkRegistrationTokenAsUsed(ctx context.Context, tokenHash string, ttl time.Duration) (alreadyUsed bool, err error)
}

type Service struct {
	txManager              database.TXManager
	usersService           usersService
	pointsService          pointsService
	notificationService    notificationService
	tokenManager           tokenManager
	refreshTokenRepository refreshTokenRepository
	registrationTokenCache registrationTokenCache
	accessTTL              time.Duration
	refreshTTL             time.Duration
	registrationTTL        time.Duration
}

func NewService(
	txManager database.TXManager,
	usersService usersService,
	pointsService pointsService,
	notificationService notificationService,
	tokenManager tokenManager,
	refreshTokenRepo refreshTokenRepository,
	registrationTokenCache registrationTokenCache,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	registrationTTL time.Duration,
) *Service {
	return &Service{
		txManager:              txManager,
		usersService:           usersService,
		pointsService:          pointsService,
		notificationService:    notificationService,
		tokenManager:           tokenManager,
		refreshTokenRepository: refreshTokenRepo,
		registrationTokenCache: registrationTokenCache,
		accessTTL:              accessTTL,
		refreshTTL:             refreshTTL,
		registrationTTL:        registrationTTL,
	}
}

func (s *Service) Login(ctx context.Context, phone string, password string) (access, refresh string, err error) {
	user, err := s.usersService.GetByPhone(ctx, phone)
	if err != nil {
		if domainErr.IsNotFound(err) {
			return "", "", domainErr.ErrInvalidCredentials
		}
		return "", "", err
	}

	if !user.Active {
		return "", "", domainErr.ErrUserInactive
	}

	valid := hash.CheckPassword(user.Password, password)
	if !valid {
		return "", "", domainErr.ErrInvalidPassword
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", err // Уже обернута в domain error
	}

	return accessToken, refreshToken, nil
}

// RegisterStaff регистрирует нового работника (orchestrator)
// Проверяет права, создает user, генерирует и отправляет ссылку для подтверждения
func (s *Service) RegisterStaff(ctx context.Context, req *command.RegisterStaffCommand) error {
	createdBy, err := session.GetSession(ctx)
	if err != nil {
		return err
	}
	// 1. Валидация прав доступа
	if validErr := s.validateStaffRegistrationPermissions(ctx, req.PointCode, req.Role, createdBy); validErr != nil {
		return validErr
	}

	// 2. Создаем user (неактивный, без пароля)
	user := &entity.User{
		Phone:       req.Phone,
		PointCode:   &req.PointCode,
		NetworkCode: ptr.Pointer(createdBy.NetworkCode()),
		Role:        req.Role,
		Active:      false, // Будет активирован после подтверждения
		UpdatedBy:   createdBy.Phone(),
	}

	if createErr := s.usersService.Create(ctx, user); createErr != nil {
		return createErr
	}

	// 3. Генерируем registration token
	token, err := s.tokenManager.NewRegistrationToken(
		&dto.RegistrationTokenPayload{
			Phone:       user.Phone,
			PointCode:   req.PointCode,
			NetworkCode: createdBy.NetworkCode(),
		}, s.registrationTTL,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to generate registration token", err)
	}

	// 4. Отправляем уведомление (WhatsApp с fallback на SMS)
	if sendErr := s.notificationService.SendRegistrationLink(ctx, req.Phone, token); sendErr != nil {
		return sendErr
	}
	return nil
}

// nolint: govet
func (s *Service) RegisterStaffConfirm(ctx context.Context, req *command.RegisterStaffConfirmCommand) (access, refresh string, err error) {
	// 1. Парсим и валидируем токен
	tokenPayload, err := s.tokenManager.ParseRegistrationToken(req.Token)
	if err != nil {
		return "", "", domainErr.NewUnauthorizedError("invalid or expired registration token").WithError(err)
	}

	// 2. Хэшируем токен для one-time use проверки
	h := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(h[:])

	var (
		user                      *entity.User
		accessToken, refreshToken string
	)

	err = s.txManager.Transaction(ctx,
		database.TXOptions{
			IsolationLevel: coreEnum.IsoLevelReadCommited,
		},
		func(txCtx context.Context) error {
			// 3. Проверяем и помечаем токен использованным (атомарно через SET NX)
			alreadyUsed, markErr := s.registrationTokenCache.MarkRegistrationTokenAsUsed(txCtx, tokenHash, 24*time.Hour)
			if markErr != nil {
				return markErr
			}
			if alreadyUsed {
				return domainErr.NewConflictError("registration token already used", nil)
			}

			// 4. Получаем пользователя
			user, err = s.usersService.GetByPhone(txCtx, tokenPayload.Phone)
			if err != nil {
				return err
			}

			// 5. Defense in depth: проверка, что пользователь еще не активирован
			if user.Active {
				return domainErr.ErrUserActivated
			}

			// 6. Хэшируем пароль
			passwordHash, err := hash.Password(req.Password)
			if err != nil {
				return domainErr.NewInternalError("failed to hash password", err)
			}

			// 7. Обновляем данные пользователя
			user.Name = req.Name
			user.Surname = req.Surname
			user.Password = passwordHash
			user.Active = true

			err = s.usersService.Update(txCtx, user)
			if err != nil {
				return err
			}
			accessToken, refreshToken, err = s.generateTokens(txCtx, user)
			return err
		},
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// validateStaffRegistrationPermissions проверяет права на регистрацию работника
func (s *Service) validateStaffRegistrationPermissions(ctx context.Context, pointCode string, role enum.Role, createdBy *session.Session) error {
	switch createdBy.Role() {
	case enum.RoleStaff:
		// Staff не может создавать пользователей
		return domainErr.NewForbiddenError("staff cannot create users")

	case enum.RoleManager:
		// Point Manager может создавать только в своей точке
		if !role.OneOf(enum.RoleStaff) {
			return domainErr.NewForbiddenError("manager cannot create (net)managers")
		}
		if pointCode != createdBy.PointCode() {
			return domainErr.NewForbiddenError("point manager can only create staff in their own point").
				WithDetail("allowed_point", createdBy.PointCode()).
				WithDetail("requested_point", pointCode)
		}

	case enum.RoleNetManager, enum.RoleSelfOwner:
		if !role.OneOf(enum.RoleManager, enum.RoleStaff) {
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

// VerifyAccessToken парсит access токен и создает domain.Session
// Это метод согласно Clean Architecture (JWT → DTO → Session)
// Service не зависит от конкретных ошибок JWT - просто оборачивает любую ошибку
func (s *Service) VerifyAccessToken(_ context.Context, tokenStr string) (*session.Session, error) {
	// 1. Парсим токен через JWT (получаем DTO из domain)
	parsed, err := s.tokenManager.ParseAccessToken(tokenStr)
	if err != nil {
		// 2. Любая ошибка парсинга → Unauthorized
		// Детали ошибки сохраняются через wrapping для логирования
		return nil, domainErr.NewUnauthorizedError("authentication failed").WithError(err)
	}

	// 3. Валидируем роль (domain logic!)
	role, err := enum.RoleString(parsed.Role)
	if err != nil {
		return nil, domainErr.ErrInvalidTokenClaims.
			WithDetail("claim", "role").
			WithError(err)
	}

	// 4. Создаем Session (domain entity)
	sess := session.NewSession().
		SetPhone(parsed.Phone).
		SetRole(role).
		SetPointCode(parsed.PointCode).
		SetNetworkCode(parsed.NetworkCode)

	return sess, nil
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
			return "", "", domainErr.ErrInvalidCredentials
		}
		return "", "", domainErr.NewInternalError("failed to get user", err)
	}

	// 5. Проверяем активность пользователя
	if !user.Active {
		return "", "", domainErr.ErrUserInactive
	}

	// 6. Генерируем новую пару токенов (автоматически сохраняет новый refresh hash)
	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) generateTokens(ctx context.Context, user *entity.User) (accessToken, refreshToken string, err error) {
	access, err := s.tokenManager.NewAccessToken(user, s.accessTTL)
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
