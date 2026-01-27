package bagsynotifications

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/notifications"
	"github.com/google/uuid"
)

// notificationSender интерфейс для отправки уведомлений
type notificationSender interface {
	SendBagsyReminder(ctx context.Context, phone string, reminder *notifications.BagsyReminder, reminderType string, recipientType string) error
}

// servicesGetter интерфейс для получения услуги
type servicesGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
}

// usersGetter интерфейс для получения пользователя (мастера)
type usersGetter interface {
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
}

// pointsGetter интерфейс для получения точки
type pointsGetter interface {
	GetByCode(ctx context.Context, code string) (*point.Point, error)
}

// MessengerAdapter адаптер для отправки уведомлений о записях
// Собирает информацию из разных сервисов и формирует сообщение
type MessengerAdapter struct {
	notificationSender notificationSender
	servicesGetter     servicesGetter
	usersGetter        usersGetter
	pointsGetter       pointsGetter
	timezone           *time.Location
}

// NewMessengerAdapter создает новый адаптер
func NewMessengerAdapter(
	notificationSender notificationSender,
	servicesGetter servicesGetter,
	usersGetter usersGetter,
	pointsGetter pointsGetter,
	timezone *time.Location,
) *MessengerAdapter {
	if timezone == nil {
		timezone = time.UTC
	}
	return &MessengerAdapter{
		notificationSender: notificationSender,
		servicesGetter:     servicesGetter,
		usersGetter:        usersGetter,
		pointsGetter:       pointsGetter,
		timezone:           timezone,
	}
}

// SendBagsyReminder собирает информацию и отправляет уведомление
func (a *MessengerAdapter) SendBagsyReminder(ctx context.Context, phone string, b *bagsy.Bagsy, notifType notification.Type, recipientType notification.RecipientType) error {
	// Получаем информацию об услуге
	serviceName := "Услуга"
	svc, err := a.servicesGetter.GetByID(ctx, b.ServiceID)
	if err == nil && svc != nil {
		serviceName = svc.Name
	}

	// Получаем информацию о мастере
	masterName := "Мастер"
	master, err := a.usersGetter.GetByPhone(ctx, b.MasterPhone)
	if err == nil && master != nil {
		masterName = master.Name
		if master.Surname != "" {
			masterName = master.Name + " " + master.Surname
		}
	}

	// Получаем информацию о точке
	pointName := "Салон"
	pt, err := a.pointsGetter.GetByCode(ctx, b.PointCode)
	if err == nil && pt != nil {
		pointName = pt.Name
	}

	// Форматируем время
	startAt := b.StartAt.In(a.timezone).Format("02.01.2006 15:04")

	reminder := &notifications.BagsyReminder{
		ServiceName: serviceName,
		MasterName:  masterName,
		PointName:   pointName,
		StartAt:     startAt,
	}

	return a.notificationSender.SendBagsyReminder(ctx, phone, reminder, notifType.String(), recipientType.String())
}
