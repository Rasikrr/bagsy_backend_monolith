package bagsies

import (
	"github.com/shopspring/decimal"
	"sort"
	"time"

	timeutil "github.com/Rasikrr/bagsy_backend_monolith/internal/util/time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

// Константа для периода слотов (4 недели)
const defaultSlotsPeriodDays = 28

type getSlotsRequest struct {
	PointCode   string  `json:"point_code" validate:"required"`
	ServiceID   string  `json:"service_id" validate:"required,uuid"`
	MasterPhone *string `json:"master_phone" validate:"omitempty,min=10"`
}

func (r *getSlotsRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *getSlotsRequest) toDomain() (*bagsy.GetAvailableSlotsCommand, error) {
	serviceID, err := uuid.Parse(r.ServiceID)
	if err != nil {
		return nil, domainErr.NewInvalidInputError("invalid service_id format", err)
	}

	// Используем UTC для всех вычислений
	now := time.Now().UTC()
	return &bagsy.GetAvailableSlotsCommand{
		PointCode:   r.PointCode,
		ServiceID:   serviceID,
		MasterPhone: r.MasterPhone,
		StartDate:   now,
		EndDate:     now.AddDate(0, 0, defaultSlotsPeriodDays),
	}, nil
}

type getSlotsResponse struct {
	ServiceID       uuid.UUID `json:"service_id"`
	PointCode       string    `json:"point_code"`
	DurationMinutes int       `json:"duration_minutes"`
	AvailableDates  []string  `json:"available_dates"`
}

func newGetSlotsResponse(slots *bagsy.AvailableSlots) *getSlotsResponse {
	// Собираем уникальные даты из всех слотов всех мастеров
	dateSet := make(map[string]struct{})

	for _, ms := range slots.MasterSlots {
		for _, ts := range ms.Slots {
			// Конвертируем в Almaty и берем только дату
			almatyTime := timeutil.ConvertUTCToAlmatyTime(ts.StartAt)
			dateStr := almatyTime.Format("2006-01-02")
			dateSet[dateStr] = struct{}{}
		}
	}

	// Конвертируем map в отсортированный slice
	availableDates := make([]string, 0, len(dateSet))
	for date := range dateSet {
		availableDates = append(availableDates, date)
	}

	// Сортируем даты
	for i := range len(availableDates) - 1 {
		for j := i + 1; j < len(availableDates); j++ {
			if availableDates[i] > availableDates[j] {
				availableDates[i], availableDates[j] = availableDates[j], availableDates[i]
			}
		}
	}

	return &getSlotsResponse{
		ServiceID:       slots.ServiceID,
		PointCode:       slots.PointCode,
		DurationMinutes: slots.DurationMinutes,
		AvailableDates:  availableDates,
	}
}

type createBagsyRequest struct {
	ServiceID   uuid.UUID `json:"service_id" validate:"required"`
	MasterPhone string    `json:"master_phone" validate:"required"`
	StartAt     time.Time `json:"start_at" validate:"required"`
	ClientPhone string    `json:"client_phone" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Surname     string    `json:"surname" validate:"required"`
	Comment     *string   `json:"comment"`
}

func (c *createBagsyRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (c *createBagsyRequest) toDomain() *bagsy.CreateBagsyCommand {
	return &bagsy.CreateBagsyCommand{
		ServiceID:   c.ServiceID,
		MasterPhone: c.MasterPhone,
		StartAt:     c.StartAt,
		ClientPhone: c.ClientPhone,
		Name:        c.Name,
		Surname:     c.Surname,
		Comment:     c.Comment,
	}
}

type createBagsyResponse struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
}

func newCreateBagsyResponse(bagsyID uuid.UUID) *createBagsyResponse {
	return &createBagsyResponse{
		BagsyID: bagsyID,
	}
}

type confirmBagsyRequest struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
	Code    string    `json:"code" validate:"required"`
}

func (c *confirmBagsyRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type resentCodeRequest struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
}

func (r *resentCodeRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

// ========== GET SLOTS FOR DAY ==========

type getSlotsForDayRequest struct {
	PointCode   string  `json:"point_code" validate:"required"`
	ServiceID   string  `json:"service_id" validate:"required,uuid"`
	Date        string  `json:"date" validate:"required"`
	MasterPhone *string `json:"master_phone" validate:"omitempty,min=10"`
}

func (r *getSlotsForDayRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	_, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return domainErr.NewValidationError("invalid date format, expected YYYY-MM-DD", err.Error())
	}
	return nil
}

func (r *getSlotsForDayRequest) toDomain() (*bagsy.GetAvailableSlotsCommand, error) {
	serviceID, err := uuid.Parse(r.ServiceID)
	if err != nil {
		return nil, domainErr.NewInvalidInputError("invalid service_id format", err)
	}

	date, _ := time.Parse("2006-01-02", r.Date)
	almatyLoc := time.FixedZone("Asia/Almaty", 5*60*60)
	startOfDayAlmaty := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, almatyLoc)
	endOfDayAlmaty := startOfDayAlmaty.Add(24 * time.Hour)

	return &bagsy.GetAvailableSlotsCommand{
		PointCode:   r.PointCode,
		ServiceID:   serviceID,
		MasterPhone: r.MasterPhone,
		StartDate:   startOfDayAlmaty.UTC(),
		EndDate:     endOfDayAlmaty.UTC(),
	}, nil
}

type masterSlotsResponse struct {
	MasterPhone        string   `json:"master_phone"`
	MasterName         string   `json:"master_name"`
	MasterServicePrice float64  `json:"master_service_price"`
	Slots              []string `json:"slots"`
}

type getSlotsForDayResponse struct {
	ServiceID       uuid.UUID             `json:"service_id"`
	PointCode       string                `json:"point_code"`
	Date            string                `json:"date"`
	DurationMinutes int                   `json:"duration_minutes"`
	Masters         []masterSlotsResponse `json:"masters"`
}

func newGetSlotsForDayResponse(slots *bagsy.AvailableSlots, date string) *getSlotsForDayResponse {
	masters := make([]masterSlotsResponse, 0, len(slots.MasterSlots))

	for _, ms := range slots.MasterSlots {
		slotTimes := make([]string, 0, len(ms.Slots))
		for _, ts := range ms.Slots {
			startAlmaty := timeutil.ConvertUTCToAlmatyTime(ts.StartAt)
			slotTimes = append(slotTimes, startAlmaty.Format("15:04"))
		}
		sort.Strings(slotTimes)

		var (
			masterServicePrice float64
		)
		if !decimal.Decimal.IsZero(ms.MasterServicePrice) {
			masterServicePrice, _ = ms.MasterServicePrice.Float64()
		}

		masters = append(masters, masterSlotsResponse{
			MasterPhone:        ms.MasterPhone,
			MasterName:         ms.MasterName,
			MasterServicePrice: masterServicePrice,
			Slots:              slotTimes,
		})
	}

	return &getSlotsForDayResponse{
		ServiceID:       slots.ServiceID,
		PointCode:       slots.PointCode,
		Date:            date,
		DurationMinutes: slots.DurationMinutes,
		Masters:         masters,
	}
}
