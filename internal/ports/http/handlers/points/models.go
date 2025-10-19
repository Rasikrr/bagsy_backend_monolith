package points

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	timeutil "github.com/Rasikrr/bagsy_backend_monolith/internal/util/time"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/deref"
)

//go:generate easyjson -all models.go
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	Coordinates Coordinates `json:"coordinates"`
	Street      string      `json:"street"`
	City        string      `json:"city"`
}

type Schedule struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open,omitempty"`
	Close   string `json:"close,omitempty"`
	AllDay  bool   `json:"all_day,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type CreatePointRequest struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	NetworkCode string     `json:"network_code,omitempty"`
	CategoryID  int        `json:"category_id,omitempty"`
	Address     Address    `json:"address"`
	City        string     `json:"city"`
	Active      bool       `json:"active"`
	Schedule    []Schedule `json:"schedule,omitempty"`
}

// UpdatePointRequest запрашивает обновление точки
type UpdatePointRequest struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	NetworkCode *string     `json:"network_code,omitempty"`
	CategoryID  *int        `json:"category_id,omitempty"`
	Address     *Address    `json:"address,omitempty"`
	City        *string     `json:"city,omitempty"`
	Active      *bool       `json:"active,omitempty"`
	Schedule    *[]Schedule `json:"schedule,omitempty"`
	UpdatedBy   string      `json:"updated_by"`
}

// PointResponse представляет точку
type PointResponse struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	NetworkCode *string    `json:"network_code,omitempty"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Address     Address    `json:"address"`
	City        string     `json:"city"`
	Active      bool       `json:"active"`
	Schedule    []Schedule `json:"schedule,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   *string    `json:"updated_at,omitempty"`
	DeletedAt   *string    `json:"deleted_at,omitempty"`
	UpdatedBy   string     `json:"updated_by"`
}

func (req CreatePointRequest) ToEntity(sessionPhone string) *entity.Point {
	return &entity.Point{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		NetworkCode: req.NetworkCode,
		CategoryID:  req.CategoryID,
		Address: entity.Address{
			Coordinates: entity.Coordinates{
				Latitude:  req.Address.Coordinates.Latitude,
				Longitude: req.Address.Coordinates.Longitude,
			},
			Street: req.Address.Street,
			City:   req.Address.City,
		},
		City:      req.City,
		Active:    req.Active,
		Schedule:  convertSchedules(req.Schedule),
		UpdatedBy: sessionPhone,
	}
}

func (req UpdatePointRequest) ToEntity(code string) *entity.Point {
	p := &entity.Point{
		Code:        code,
		Name:        deref.String(req.Name),
		Description: req.Description,
		NetworkCode: deref.String(req.NetworkCode),
		CategoryID:  deref.Int(req.CategoryID),
		City:        deref.String(req.City),
		Active:      deref.Bool(req.Active),
		UpdatedBy:   req.UpdatedBy,
	}
	if req.Address != nil {
		p.Address = entity.Address{
			Coordinates: entity.Coordinates{
				Latitude:  req.Address.Coordinates.Latitude,
				Longitude: req.Address.Coordinates.Longitude,
			},
			Street: req.Address.Street,
			City:   req.Address.City,
		}
	}
	if req.Schedule != nil {
		p.Schedule = convertSchedules(*req.Schedule)
	}
	return p
}

func convertSchedules(s []Schedule) []entity.Schedule {
	result := make([]entity.Schedule, len(s))

	for i, v := range s {
		result[i] = entity.Schedule{
			Open:    timeutil.ConvertStrToScheduleTime(v.Open),
			Close:   timeutil.ConvertStrToScheduleTime(v.Close),
			AllDay:  v.AllDay,
			Comment: v.Comment,
		}
	}

	return result
}
