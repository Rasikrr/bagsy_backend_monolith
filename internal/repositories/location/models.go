package location

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	ID                  uuid.UUID  `db:"id"`
	OrganizationID      uuid.UUID  `db:"organization_id"`
	CategoryID          uuid.UUID  `db:"category_id"`
	Name                string     `db:"name"`
	Description         *string    `db:"description"`
	Phone               *string    `db:"phone"`
	Slug                string     `db:"slug"`
	City                *string    `db:"city"`
	Street              *string    `db:"address_street"`
	Building            *string    `db:"address_building"`
	Details             *string    `db:"address_details"`
	Longitude           *float64   `db:"longitude"`
	Latitude            *float64   `db:"latitude"`
	Active              bool       `db:"active"`
	ScheduleType        string     `db:"schedule_type"`
	SlotDurationMinutes int        `db:"slot_duration_minutes"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at"`
}

func fromDomain(l *location.Location) *model {
	var phone *string
	if l.Phone != nil {
		p := l.Phone.String()
		phone = &p
	}

	var city, street, building, details *string
	if l.Address != nil {
		city = &l.Address.City
		street = &l.Address.Street
		building = &l.Address.Building
		details = &l.Address.Details
	}

	var lat, lng *float64
	if l.Coordinates != nil {
		lat = &l.Coordinates.Latitude
		lng = &l.Coordinates.Longitude
	}

	return &model{
		ID:                  l.ID,
		OrganizationID:      l.OrganizationID,
		CategoryID:          l.CategoryID,
		Name:                l.Name,
		Description:         l.Description,
		Phone:               phone,
		Slug:                l.Slug.String(),
		City:                city,
		Street:              street,
		Building:            building,
		Details:             details,
		Longitude:           lng,
		Latitude:            lat,
		Active:              l.Active,
		ScheduleType:        string(l.ScheduleType),
		SlotDurationMinutes: l.SlotDurationMinutes.Minutes(),
		CreatedAt:           l.CreatedAt,
		UpdatedAt:           l.UpdatedAt,
		DeletedAt:           l.DeletedAt,
	}
}

func (m *model) toDomain() (*location.Location, error) {
	var phone *shared.Phone
	if m.Phone != nil {
		p, err := shared.NewPhone(*m.Phone)
		if err != nil {
			return nil, err
		}
		phone = &p
	}

	var addr *location.Address
	if m.City != nil {
		addr = &location.Address{
			City:     *m.City,
			Street:   *getStringValue(m.Street),
			Building: *getStringValue(m.Building),
			Details:  *getStringValue(m.Details),
		}
	}

	var coords *location.Coordinates
	if m.Latitude != nil && m.Longitude != nil {
		coords = &location.Coordinates{
			Latitude:  *m.Latitude,
			Longitude: *m.Longitude,
		}
	}
	slug, err := shared.NewSlug(m.Slug)
	if err != nil {
		return nil, err
	}
	dur, err := shared.NewDuration(m.SlotDurationMinutes)
	if err != nil {
		return nil, err
	}

	return &location.Location{
		ID:                  m.ID,
		OrganizationID:      m.OrganizationID,
		CategoryID:          m.CategoryID,
		Name:                m.Name,
		Description:         m.Description,
		Phone:               phone,
		Slug:                slug,
		Address:             addr,
		Coordinates:         coords,
		Active:              m.Active,
		ScheduleType:        location.ScheduleType(m.ScheduleType),
		SlotDurationMinutes: dur,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		DeletedAt:           m.DeletedAt,
	}, nil
}

func getStringValue(s *string) *string {
	if s == nil {
		v := ""
		return &v
	}
	return s
}
