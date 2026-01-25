package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type getUsersRequest struct {
	PointCode   *string  `query:"point_code" validate:""`
	NetworkCode *string  `query:"network_code"`
	Roles       []string `query:"role"`
	PhoneSearch *string  `query:"phone_search"` // Частичный или полный поиск по номеру телефона
	Limit       uint64   `query:"limit" validate:"max=100"`
	Offset      uint64   `query:"offset" validate:"min=0"`
	OrderBy     string   `query:"order_by" validate:"oneof=phone name surname point_code network_code created_at updated_at"`
	SortOrder   string   `query:"sort_order" validate:"oneof=asc desc"`
}

func (r *getUsersRequest) GetQueryParameters(req *http.Request) error {
	q := req.URL.Query()

	if pointCode := q.Get("point_code"); pointCode != "" {
		r.PointCode = &pointCode
	}

	if networkCode := q.Get("network_code"); networkCode != "" {
		r.NetworkCode = &networkCode
	}

	if roles := q["role"]; len(roles) > 0 {
		r.Roles = roles
	}

	if phoneSearch := q.Get("phone_search"); phoneSearch != "" {
		r.PhoneSearch = &phoneSearch
	}

	r.Limit = 20
	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			r.Limit = uint64(limit)
		}
	}

	r.Offset = 0
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset > 0 {
			r.Offset = uint64(offset)
		}
	}

	r.OrderBy = "created_at"
	if sortBy := q.Get("sort_by"); sortBy != "" {
		r.OrderBy = sortBy
	}

	r.SortOrder = "asc"
	if sortOrder := q.Get("sort_order"); sortOrder != "" {
		r.SortOrder = sortOrder
	}

	return nil
}

func (r *getUsersRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *getUsersRequest) toFilter() (*user.Filter, error) {
	q := &user.Filter{
		PointCode:   r.PointCode,
		NetworkCode: r.NetworkCode,
		PhoneSearch: r.PhoneSearch,
		Limit:       r.Limit,
		Offset:      r.Offset,
		OrderBy:     r.OrderBy,
	}
	sortOrder, err := enum.SortOrderString(r.SortOrder)
	if err != nil {
		return nil, domainErr.NewValidationError("sort_order contains an invalid value").
			WithDetail("sort_order", r.SortOrder)
	}
	q.SortOrder = sortOrder

	roles := make([]user.Role, len(r.Roles))
	for i, role := range r.Roles {
		roleEnum, enumErr := user.RoleString(role)
		if enumErr != nil {
			return nil, domainErr.NewValidationError("roles contains invalid item").
				WithDetail("role", role)
		}
		roles[i] = roleEnum
	}
	q.Roles = roles

	return q, nil
}

type staffScheduleDTO struct {
	WeekDay int    `json:"week_day"`
	Open    time.Time `json:"open"`
	Close   time.Time `json:"close"`
	AllDay  bool   `json:"all_day"`
	Comment string `json:"comment"`
}

type userDTO struct {
	Phone       string              `json:"phone"`
	Role        string              `json:"role"`
	Name        string              `json:"name"`
	Surname     string              `json:"surname"`
	PointCode   *string             `json:"point_code,omitempty"`
	NetworkCode *string             `json:"network_code,omitempty"`
	Active      bool                `json:"active"`
	AvatarURL   string              `json:"avatar_url,omitempty"`
	Schedule    []*staffScheduleDTO `json:"schedule,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
}

func toStaffScheduleDTO(schedule *user.ScheduleElement) *staffScheduleDTO {
	return &staffScheduleDTO{
		WeekDay: schedule.WeekDay,
		Open:    schedule.Open,
		Close:   schedule.Close,
		AllDay:  schedule.AllDay,
		Comment: schedule.Comment,
	}
}

func toUserDTO(user *user.User) *userDTO {
	d := &userDTO{
		Phone:       user.Phone,
		Role:        user.Role.String(),
		Name:        user.Name,
		Surname:     user.Surname,
		PointCode:   user.PointCode,
		NetworkCode: user.NetworkCode,
		Active:      user.Active,
		CreatedAt:   user.CreatedAt,
	}

	if user.UpdatedAt != nil {
		updatedAt := user.UpdatedAt
		d.UpdatedAt = updatedAt
	}

	if len(user.Schedule) > 0 {
		schedules := make([]*staffScheduleDTO, 0, len(user.Schedule))
		for _, s := range user.Schedule {
			schedules = append(schedules, toStaffScheduleDTO(s))
		}
		d.Schedule = schedules
	}
	if user.Avatar != nil {
		d.AvatarURL = user.Avatar.GetURL()
	}

	return d
}

type getUsersResponse struct {
	Users []*userDTO `json:"users"`
	Total int        `json:"total"`
}

func toGetUsersResponse(paginated *query.Page[*user.User]) getUsersResponse {
	dtos := make([]*userDTO, 0, len(paginated.Items))
	for _, user := range paginated.Items {
		dtos = append(dtos, toUserDTO(user))
	}

	return getUsersResponse{
		Users: dtos,
		Total: paginated.Total,
	}
}

type updateUserRequest struct {
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	AvatarID *uuid.UUID `json:"avatar_id"`
}

func (r *updateUserRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *updateUserRequest) toDomain() *user.UpdateUserCommand {
	return &user.UpdateUserCommand{
		Name:     r.Name,
		Surname:  r.Surname,
		AvatarID: r.AvatarID,
	}
}

type scheduleRequestDTO struct {
	WeekDay int       `json:"week_day" validate:"min=0,max=6"`
	From    time.Time `json:"from" validate:"required"`
	To      time.Time `json:"to" validate:"required"`
	AllDay  bool      `json:"all_day"`
	Comment string    `json:"comment"`
}

type updateScheduleRequest struct {
	Schedule []scheduleRequestDTO `json:"schedule" validate:"required"`
}

func (r *updateScheduleRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}

	if len(r.Schedule) < 1 {
		return domainErr.NewValidationError("schedule must contain exactly 7 days").
			WithDetail("length", len(r.Schedule))
	}

	return nil
}

func (r *updateScheduleRequest) toDomain() (user.Schedule, error) {
	schedules := make(user.Schedule, 0, len(r.Schedule))

	for _, s := range r.Schedule {
		schedules = append(schedules, &user.ScheduleElement{
			WeekDay: s.WeekDay,
			Open:    s.From.UTC(),
			Close:   s.To.UTC(),
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}
	return schedules, nil
}
