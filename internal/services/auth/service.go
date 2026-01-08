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
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
)

const (
	maxVerificationAttempts = 3
)

type usersService interface {
	Create(ctx context.Context, user *entity.User, password string) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	UpdateWithPassword(ctx context.Context, user *entity.User, rawPassword string) error
	UpdatePasswordByPhone(ctx context.Context, phone, rawPassword string) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
}

type networkService interface {
	CreateForRegistration(ctx context.Context, req *command.CreateNetworkCommand, createdBy string) (*entity.Network, error)
}

type notificationService interface {
	SendStaffRegistrationLink(ctx context.Context, phone, token string) error
	SendManagementAuthConfirmationCode(ctx context.Context, phone, code string) error
	SendPasswordChangeLink(ctx context.Context, phone, token string) error
}

type tokenManager interface {
	NewAccessToken(user *entity.User, ttl time.Duration) (string, error)
	NewRefreshToken() (raw, hash string, err error)
	NewAuthToken(dto *dto.RegistrationTokenPayload, ttl time.Duration) (string, error)
	ParseAccessToken(token string) (*dto.AccessTokenPayload, error)
	ParseAuthToken(token string) (*dto.RegistrationTokenPayload, error)
}

type registerCache interface {
	SaveManagementRequest(ctx context.Context, req *command.RegisterManagementCommand) error
	GetManagementRequest(ctx context.Context, phone string) (*command.RegisterManagementCommand, error)
	DeleteManagementRequest(ctx context.Context, phone string) error

	SaveStaffRequest(ctx context.Context, req *command.RegisterStaffCommand) error
	GetStaffRequest(ctx context.Context, phone string) (*command.RegisterStaffCommand, error)
	DeleteStaffRequest(ctx context.Context, phone string) error
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
type tokensCache interface {
	// MarkRegistrationTokenAsUsed помечает токен как использованный (атомарно через SET NX)
	// Возвращает true если токен уже был использован ранее
	MarkRegistrationTokenAsUsed(ctx context.Context, tokenHash string, ttl time.Duration) (alreadyUsed bool, err error)
}

type Service struct {
	txManager              database.TXManager
	usersService           usersService
	pointsService          pointsService
	notificationService    notificationService
	networkService         networkService
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
	usersService usersService,
	pointsService pointsService,
	notificationService notificationService,
	networkService networkService,
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
		usersService:           usersService,
		pointsService:          pointsService,
		notificationService:    notificationService,
		networkService:         networkService,
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

func (s *Service) RegisterManagement(ctx context.Context, req *command.RegisterManagementCommand) error {
	authCode := codegen.GenerateAuthCode()

	req.AuthCode = authCode

	// Проверка, что запрос не был уже создан ранее
	existedReq, err := s.registerCache.GetManagementRequest(ctx, req.Phone)
	if err != nil {
		if !domainErr.IsNotFound(err) {
			return err
		}
	}
	if existedReq != nil {
		return domainErr.NewConflictError("registration request already created", nil)
	}

	err = s.registerCache.SaveManagementRequest(ctx, req)
	if err != nil {
		return err
	}
	return s.notificationService.SendManagementAuthConfirmationCode(ctx, req.Phone, authCode)
}

func (s *Service) ResendRegisterManagementCode(ctx context.Context, phone string) error {
	req, err := s.registerCache.GetManagementRequest(ctx, phone)
	if err != nil {
		return err
	}
	newCode := codegen.GenerateAuthCode()
	req.AuthCode = newCode
	req.Attempts = 0
	err = s.registerCache.SaveManagementRequest(ctx, req)
	if err != nil {
		return err
	}
	return s.notificationService.SendManagementAuthConfirmationCode(ctx, phone, newCode)
}

func (s *Service) RegisterManagementConfirm(ctx context.Context, phone, code string) (access, refresh string, err error) {
	req, err := s.registerCache.GetManagementRequest(ctx, phone)
	if err != nil {
		return "", "", err
	}
	if req.AuthCode != code {
		req.Attempts++
		if req.Attempts >= maxVerificationAttempts {
			// Удаляем данные регистрации после превышения лимита попыток
			if delErr := s.registerCache.DeleteManagementRequest(ctx, phone); delErr != nil {
				return "", "", delErr
			}
			return "", "", domainErr.ErrTooManyVerificationAttempts
		}
		// Сохраняем обновленный счетчик попыток
		if saveErr := s.registerCache.SaveManagementRequest(ctx, req); saveErr != nil {
			return "", "", saveErr
		}
		return "", "", domainErr.ErrInvalidVerificationCode.
			WithDetail("attempts_remaining", maxVerificationAttempts-req.Attempts)
	}
	var user *entity.User

	err = s.txManager.Transaction(ctx, database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}, func(ctx context.Context) error {
		network, netErr := s.networkService.CreateForRegistration(ctx, &command.CreateNetworkCommand{
			Name:        req.NetworkRegisterInfo.Name,
			Description: req.NetworkRegisterInfo.Description,
		}, req.Phone)
		if netErr != nil {
			return netErr
		}
		user, err = s.usersService.GetByPhone(ctx, req.Phone)
		if err != nil {
			if !domainErr.IsNotFound(err) {
				return err
			}
			// Надо создать
			user = &entity.User{
				Name:        req.Name,
				Surname:     req.Surname,
				Phone:       req.Phone,
				Role:        req.Role,
				NetworkCode: &network.Code,
				Active:      true,
			}
			err = s.usersService.Create(ctx, user, req.Password)
		} else {
			user.Name = req.Name
			user.Surname = req.Surname
			user.Role = req.Role
			user.NetworkCode = &network.Code
			user.Active = true
			err = s.usersService.UpdateWithPassword(ctx, user, req.Password)
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return s.generateTokens(ctx, user)
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
	req.NetworkCode = createdBy.NetworkCode()

	// 2. Проверим что запрос не был создан ранее
	existedReq, err := s.registerCache.GetStaffRequest(ctx, req.Phone)
	if err != nil {
		if !domainErr.IsNotFound(err) {
			return err
		}
	}
	if existedReq != nil {
		return domainErr.NewConflictError("registration request already created", nil)
	}

	// 3. Сохраним данные в кеше чтобы преждевременно не создавать юзера
	err = s.registerCache.SaveStaffRequest(ctx, req)
	if err != nil {
		return err
	}
	// 4. Генерируем registration token
	token, err := s.tokenManager.NewAuthToken(
		&dto.RegistrationTokenPayload{
			Phone: req.Phone,
		}, s.registrationTTL,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to generate registration token", err)
	}

	// 3. Отправляем уведомление (WhatsApp с fallback на SMS)
	if sendErr := s.notificationService.SendStaffRegistrationLink(ctx, req.Phone, token); sendErr != nil {
		return sendErr
	}
	return nil
}

func (s *Service) ResendRegisterStaffLink(ctx context.Context, phone string) error {
	req, err := s.registerCache.GetStaffRequest(ctx, phone)
	if err != nil {
		return err
	}
	token, tErr := s.tokenManager.NewAuthToken(
		&dto.RegistrationTokenPayload{
			Phone: req.Phone,
		}, s.registrationTTL,
	)
	if tErr != nil {
		return domainErr.NewInternalError("failed to generate registration token", tErr)
	}
	if sendErr := s.notificationService.SendStaffRegistrationLink(ctx, req.Phone, token); sendErr != nil {
		return sendErr
	}
	return nil
}

// nolint: govet
func (s *Service) RegisterStaffConfirm(ctx context.Context, req *command.RegisterStaffConfirmCommand) (access, refresh string, err error) {
	// 1. Парсим и валидируем токен
	tokenPayload, err := s.tokenManager.ParseAuthToken(req.Token)
	if err != nil {
		return "", "", domainErr.NewUnauthorizedError("invalid or expired registration token").WithError(err)
	}

	createReq, err := s.registerCache.GetStaffRequest(ctx, tokenPayload.Phone)
	if err != nil {
		return "", "", err
	}

	// 2. Хэшируем токен для one-time use проверки
	h := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(h[:])

	var (
		accessToken, refreshToken string
	)

	err = s.txManager.Transaction(ctx,
		database.TXOptions{
			IsolationLevel: coreEnum.IsoLevelReadCommited,
		},
		func(txCtx context.Context) error {
			// 3. Проверяем и помечаем токен использованным (атомарно через SET NX)
			alreadyUsed, markErr := s.tokensCache.MarkRegistrationTokenAsUsed(txCtx, tokenHash, 24*time.Hour)
			if markErr != nil {
				return markErr
			}
			if alreadyUsed {
				return domainErr.NewConflictError("registration token already used", nil)
			}

			user, err := s.usersService.GetByPhone(ctx, createReq.Phone)
			if err != nil {
				if !domainErr.IsNotFound(err) {
					return err
				}
				// Значит надо создать нового юзера
				user = &entity.User{
					Phone:       createReq.Phone,
					Name:        createReq.Name,
					Surname:     createReq.Surname,
					PointCode:   &createReq.PointCode,
					NetworkCode: ptr.Pointer(createReq.NetworkCode),
					Role:        createReq.Role,
					Active:      true,
				}
				err = s.usersService.Create(ctx, user, req.Password)
			} else {
				user.Name = createReq.Name
				user.Surname = createReq.Surname
				user.PointCode = &createReq.PointCode
				user.NetworkCode = &createReq.NetworkCode
				user.Role = createReq.Role
				user.Active = true
				err = s.usersService.UpdateWithPassword(ctx, user, req.Password)
			}
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

func (s *Service) SendPasswordChangeLink(ctx context.Context, phone string) error {
	exist, err := s.usersService.ExistsByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if !exist {
		return domainErr.ErrUserNotFound
	}
	token, err := s.tokenManager.NewAuthToken(&dto.RegistrationTokenPayload{
		Phone: phone,
	}, s.registrationTTL)
	if err != nil {
		return domainErr.NewInternalError("failed to generate change password token", err)
	}
	return s.notificationService.SendPasswordChangeLink(ctx, phone, token)
}

func (s *Service) ChangePassword(ctx context.Context, req *command.ChangePasswordConfirmCommand) error {
	// 1. Парсим и валидируем токен
	tokenPayload, err := s.tokenManager.ParseAuthToken(req.Token)
	if err != nil {
		return domainErr.NewUnauthorizedError("invalid or expired token").WithError(err)
	}

	// 2. Хэшируем токен для one-time use проверки
	h := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(h[:])

	err = s.txManager.Transaction(ctx,
		database.TXOptions{
			IsolationLevel: coreEnum.IsoLevelReadCommited,
		},
		func(txCtx context.Context) error {
			// 3. Проверяем и помечаем токен использованным (атомарно через SET NX)
			alreadyUsed, markErr := s.tokensCache.MarkRegistrationTokenAsUsed(txCtx, tokenHash, 24*time.Hour)
			if markErr != nil {
				return markErr
			}
			if alreadyUsed {
				return domainErr.NewConflictError("change password token already used", nil)
			}
			return s.usersService.UpdatePasswordByPhone(ctx, tokenPayload.Phone, req.Password)
		},
	)
	return err
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
