// nolint: unused
package calendar

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type getCalendarRequest struct {
	From        time.Time `query:"from" validate:"required"`
	To          time.Time `query:"to" validate:"required"`
	PointCode   *string   `query:"point_code"`
	MasterPhone *string   `query:"master_phone"`
}

func (g *getCalendarRequest) Validate() error {
	if err := request.GetValidator().Struct(g); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (g *getCalendarRequest) GetQueryParameters(r *http.Request) error {
	pointCode := r.URL.Query().Get("point_code")
	if pointCode != "" {
		g.PointCode = &pointCode
	}
	masterPhone := r.URL.Query().Get("master_phone")
	if masterPhone != "" {
		g.MasterPhone = &masterPhone
	}

	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			return request.HandleValidationError(err)
		}
		g.From = from
	}

	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			return request.HandleValidationError(err)
		}
		g.To = to
	}
	return nil
}

func (g *getCalendarRequest) toQuery() *bagsies.GetCalendarQuery {
	return &bagsies.GetCalendarQuery{
		StartDate:   g.From,
		EndDate:     g.To,
		PointCode:   g.PointCode,
		MasterPhone: g.MasterPhone,
	}
}

type bagsyInfoDTO struct {
	ID          uuid.UUID  `json:"id"`
	PointCode   string     `json:"point_code"`
	ClientPhone string     `json:"client_phone"`
	MasterPhone string     `json:"master_phone"`
	Status      string     `json:"status"`
	Price       float64    `json:"price"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	Comment     *string    `json:"comment,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func newBagsyInfoDTO(b *bagsy.Bagsy) *bagsyInfoDTO {
	return &bagsyInfoDTO{
		ID:          b.ID,
		PointCode:   b.PointCode,
		ClientPhone: b.ClientPhone,
		MasterPhone: b.MasterPhone,
		Status:      b.Status.String(),
		Price:       b.Price.InexactFloat64(),
		StartAt:     b.StartAt,
		EndAt:       b.EndAt,
		Comment:     b.Comment,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

type serviceCategoryDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func newCategoryDTO(c service.Category) serviceCategoryDTO {
	return serviceCategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type serviceSubcategoryDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func newSubcategoryDTO(s *service.Subcategory) *serviceSubcategoryDTO {
	if s == nil {
		return nil
	}
	return &serviceSubcategoryDTO{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

type serviceInfoDTO struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	DurationMinutes int        `json:"duration_minutes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	Color           string     `json:"color"`
}

func newServiceInfoDTO(s service.Service) *serviceInfoDTO {
	return &serviceInfoDTO{
		ID:              s.ID,
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		Color:           s.Color.String(),
	}
}

type calendarDTO struct {
	BagsyInfo   *bagsyInfoDTO   `json:"bagsy_info"`
	ServiceInfo *serviceInfoDTO `json:"service_info"`
}

func newCalendarDTO(b *bagsies.CalendarElement) *calendarDTO {
	return &calendarDTO{
		BagsyInfo:   newBagsyInfoDTO(b.Bagsy),
		ServiceInfo: newServiceInfoDTO(*b.Service),
	}
}

type calendarResponseDTO struct {
	Calendar []*calendarDTO `json:"calendar"`
}

func newCalendarResponse(list []*bagsies.CalendarElement) *calendarResponseDTO {
	out := make([]*calendarDTO, 0, len(list))
	for _, el := range list {
		out = append(out, newCalendarDTO(el))
	}
	return &calendarResponseDTO{
		Calendar: out,
	}
}
