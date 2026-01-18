package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/cockroachdb/errors"
)

type GetCalendarQuery struct {
	StartDate   time.Time
	EndDate     time.Time
	PointCode   *string
	MasterPhone *string
}

type CalendarElement struct {
	Bagsy   *bagsy.Bagsy
	Service *service.Service
}

func newCalendarElement(bagsy *bagsy.Bagsy, service *service.Service) (*CalendarElement, error) {
	if bagsy == nil || service == nil {
		return nil, errors.New("bagsy or service is nil")
	}
	return &CalendarElement{
		Bagsy:   bagsy,
		Service: service,
	}, nil
}
