package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util"
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
	SendWhatsApp(ctx context.Context, phone, message string) error
	SendSMS(ctx context.Context, phone, message string) error
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
	registrationConfirmURL string
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
	registrationConfirmURL string,
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
		registrationConfirmURL: registrationConfirmURL,
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
func (s *Service) RegisterStaff(ctx context.Context, phone, pointCode string) error {
	createdBy, err := session.GetSession(ctx)
	if err != nil {
		return err
	}
	// 1. Валидация прав доступа
	if err := s.validateStaffRegistrationPermissions(ctx, pointCode, createdBy); err != nil {
		return err
	}

	// 2. Создаем user (неактивный, без пароля)
	user := &entity.User{
		Phone:       phone,
		PointCode:   &pointCode,
		NetworkCode: util.Pointer(createdBy.NetworkCode()),
		Role:        enum.RoleStaff,
		Active:      false, // Будет активирован после подтверждения
		UpdatedBy:   createdBy.Phone(),
	}

	if err := s.usersService.Create(ctx, user); err != nil {
		return err
	}

	// 3. Генерируем registration token
	token, err := s.tokenManager.NewRegistrationToken(
		&dto.RegistrationTokenPayload{
			Phone:       user.Phone,
			PointCode:   pointCode,
			NetworkCode: createdBy.NetworkCode(),
		}, s.registrationTTL,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to generate registration token", err)
	}

	// 4. Формируем ссылку для подтверждения
	link := fmt.Sprintf("%s?token=%s", s.registrationConfirmURL, token)

	// 5. Отправляем уведомление (WhatsApp с fallback на SMS)
	message := fmt.Sprintf("Добро пожаловать в Bagsy! Завершите регистрацию по ссылке: %s", link)

	if err := s.notificationService.SendWhatsApp(ctx, phone, message); err != nil {
		// Fallback на SMS если WhatsApp недоступен
		if err := s.notificationService.SendSMS(ctx, phone, message); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) RegisterStaffConfirm(ctx context.Context, req *command.RegisterStaffConfirmRequest) (access, refresh string, err error) {
	// 1. Парсим и валидируем токен
	tokenPayload, err := s.tokenManager.ParseRegistrationToken(req.Token)
	if err != nil {
		return "", "", domainErr.NewUnauthorizedError("invalid or expired registration token").WithError(err)
	}

	// 2. Хэшируем токен для one-time use проверки
	h := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(h[:])

	var user *entity.User

	err = s.txManager.Transaction(ctx,
		database.TXOptions{
			IsolationLevel: coreEnum.IsoLevelReadCommited,
		},
		func(txCtx context.Context) error {
			// 3. Проверяем и помечаем токен использованным (атомарно через SET NX)
			alreadyUsed, err := s.registrationTokenCache.MarkRegistrationTokenAsUsed(txCtx, tokenHash, 24*time.Hour)
			if err != nil {
				return err
			}
			if alreadyUsed {
				return domainErr.NewConflictError("registration token already used", nil)
			}

			// 4. Получаем пользователя
			user, err = s.usersService.GetByPhone(txCtx, tokenPayload.Phone)
			if err != nil {
				return err
			}

			// 5. Defense in depth: проверка что пользователь еще не активирован
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

			return s.usersService.Update(txCtx, user)
		},
	)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, user)
}

// validateStaffRegistrationPermissions проверяет права на регистрацию работника
func (s *Service) validateStaffRegistrationPermissions(ctx context.Context, pointCode string, createdBy *session.Session) error {
	switch createdBy.Role() {
	case enum.RoleStaff:
		// Staff не может создавать пользователей
		return domainErr.NewForbiddenError("staff cannot create users")

	case enum.RoleManager:
		// Point Manager может создавать только в своей точке
		if pointCode != createdBy.PointCode() {
			return domainErr.NewForbiddenError("point manager can only create staff in their own point").
				WithDetail("allowed_point", createdBy.PointCode()).
				WithDetail("requested_point", pointCode)
		}

	case enum.RoleNetManager, enum.RoleSelfOwner:
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
func (s *Service) VerifyAccessToken(ctx context.Context, tokenStr string) (*session.Session, error) {
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

// RefreshTokens обновляет пару токенов по refresh токену
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
	if err := s.refreshTokenRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return "", "", domainErr.NewInternalError("failed to delete old refresh token", err)
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
	if err := s.refreshTokenRepository.SaveRefreshToken(ctx, refreshHash, user.Phone, s.refreshTTL); err != nil {
		return "", "", domainErr.NewInternalError("failed to save refresh token", err)
	}

	return access, refresh, nil
}
