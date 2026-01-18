package bagsies

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// GetCalendar возвращает календарь записей
// Для Staff - только свои записи
// Для Manager+ - записи всей точки (опционально с фильтром по мастеру)
func (s *Service) GetCalendar(
	ctx context.Context,
	query *GetCalendarQuery,
) ([]*CalendarElement, error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	filter, buildErr := s.buildCalendarFilter(ctx, act, query)
	if buildErr != nil {
		return nil, buildErr
	}
	log.Infof(ctx, "GetCalendar: %+v", filter)

	bagsies, err := s.bagsiesRepository.GetOccupiedSlots(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(bagsies) == 0 {
		return []*CalendarElement{}, nil
	}

	servicesIDs := lo.Map(bagsies, func(item *bagsy.Bagsy, _ int) uuid.UUID {
		return item.ServiceID
	})

	services, err := s.servicesService.GetByIDs(ctx, servicesIDs)
	if err != nil {
		return nil, err
	}

	return s.createCalendarElements(bagsies, services)
}

// nolint: nestif
// buildCalendarFilter формирует фильтр в зависимости от роли
func (s *Service) buildCalendarFilter(
	ctx context.Context,
	act *actor.Actor,
	query *GetCalendarQuery,
) (*bagsy.OccupiedSlotsFilter, error) {
	filter := &bagsy.OccupiedSlotsFilter{
		StartAt:   query.StartDate,
		EndAt:     query.EndDate,
		PointCode: act.PointCode(),
	}
	if act.Role() == user.RoleStaff {
		filter.MasterPhones = []string{act.Phone()}
		return filter, nil
	}

	if query.PointCode != nil {
		if act.Role() == user.RoleManager {
			if *query.PointCode != act.PointCode() {
				return nil, domainErr.NewForbiddenError("forbidden to manager get bagsies from other points").
					WithDetail("queried point_code", *query.PointCode).
					WithDetail("actual point_code", act.PointCode())
			}
		}
		point, err := s.pointsService.GetByCode(ctx, *query.PointCode)
		if err != nil {
			return nil, err
		}
		if point.NetworkCode != act.NetworkCode() {
			return nil, domainErr.NewForbiddenError("forbidden to get bagsies from other networks")
		}
		filter.PointCode = point.Code
	}
	if query.MasterPhone != nil {
		filter.MasterPhones = []string{*query.MasterPhone}
	}
	return filter, nil
}

func (s *Service) createCalendarElements(bagsies []*bagsy.Bagsy, services []*service.Service) ([]*CalendarElement, error) {
	servicesMap := lo.SliceToMap(services, func(item *service.Service) (uuid.UUID, *service.Service) {
		return item.ID, item
	})

	calendarElements := make([]*CalendarElement, 0, len(bagsies))
	for _, bagsy := range bagsies {
		calendarElement, err := newCalendarElement(bagsy, servicesMap[bagsy.ServiceID])
		if err != nil {
			return nil, domainErr.NewInternalError("failed to create calendar element", err)
		}
		calendarElements = append(calendarElements, calendarElement)
	}
	return calendarElements, nil
}
