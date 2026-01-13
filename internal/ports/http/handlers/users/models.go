package users

import (
	"net/http"
	"strconv"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	timeutil "github.com/Rasikrr/bagsy_backend_monolith/internal/util/time"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type getUsersRequest struct {
	PointCode   *string  `query:"point_code" validate:""`
	NetworkCode *string  `query:"network_code"`
	Roles       []string `query:"role"`
	Phones      []string `query:"phone"`
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

	if phones := q["phone"]; len(phones) > 0 {
		r.Phones = phones
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

func (r *getUsersRequest) toFilter() (*query.UserFilter, error) {
	q := query.UserFilter{
		PointCode:   r.PointCode,
		NetworkCode: r.NetworkCode,
		Phones:      r.Phones,
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

	roles := make([]enum.Role, len(r.Roles))
	for i, role := range r.Roles {
		roleEnum, enumErr := enum.RoleString(role)
		if enumErr != nil {
			return nil, domainErr.NewValidationError("roles contains invalid item").
				WithDetail("role", role)
		}
		roles[i] = roleEnum
	}
	q.Roles = roles

	return &q, nil
}

type staffScheduleDTO struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open"`
	Close   string `json:"close"`
	AllDay  bool   `json:"all_day"`
	Comment string `json:"comment"`
}

type userDTO struct {
	Phone       string             `json:"phone"`
	Role        string             `json:"role"`
	Name        string             `json:"name"`
	Surname     string             `json:"surname"`
	PointCode   *string            `json:"point_code,omitempty"`
	NetworkCode *string            `json:"network_code,omitempty"`
	Active      bool               `json:"active"`
	AvatarURL   *string            `json:"avatar_url,omitempty"`
	Schedule    []staffScheduleDTO `json:"schedule,omitempty"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   *string            `json:"updated_at,omitempty"`
}

func toStaffScheduleDTO(schedule entity.StaffSchedule) staffScheduleDTO {
	return staffScheduleDTO{
		WeekDay: schedule.WeekDay,
		Open:    schedule.Open.Format("15:04:05"),
		Close:   schedule.Close.Format("15:04:05"),
		AllDay:  schedule.AllDay,
		Comment: schedule.Comment,
	}
}

func toUserWithAvatar(u *dto.UserWithAvatar) *userDTO {
	userDTO := toUserDTO(u.User)
	if u.AvatarURL != nil {
		userDTO.AvatarURL = u.AvatarURL
	}
	return userDTO
}

func toUserDTO(user *entity.User) *userDTO {
	d := &userDTO{
		Phone:       user.Phone,
		Role:        user.Role.String(),
		Name:        user.Name,
		Surname:     user.Surname,
		PointCode:   user.PointCode,
		NetworkCode: user.NetworkCode,
		Active:      user.Active,
		CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.UpdatedAt != nil {
		updatedAt := user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		d.UpdatedAt = &updatedAt
	}

	if len(user.Schedule) > 0 {
		schedules := make([]staffScheduleDTO, 0, len(user.Schedule))
		for _, s := range user.Schedule {
			schedules = append(schedules, toStaffScheduleDTO(s))
		}
		d.Schedule = schedules
	}

	return d
}

type getUsersResponse struct {
	Users []*userDTO `json:"users"`
	Total int        `json:"total"`
}

func toGetUsersResponse(paginated *dto.PaginatedUsers) getUsersResponse {
	dtos := make([]*userDTO, 0, len(paginated.Users))
	for _, user := range paginated.Users {
		dtos = append(dtos, toUserWithAvatar(user))
	}

	return getUsersResponse{
		Users: dtos,
		Total: paginated.Total,
	}
}

type updateUserRequest struct {
	Name     string     `json:"name" validate:"required"`
	Surname  string     `json:"surname" validate:"required"`
	AvatarID *uuid.UUID `json:"avatar_id"`
}

func (r *updateUserRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *updateUserRequest) toDomain() *command.UpdateUserCommand {
	return &command.UpdateUserCommand{
		Name:     r.Name,
		Surname:  r.Surname,
		AvatarID: r.AvatarID,
	}
}

type scheduleRequestDTO struct {
	WeekDay int    `json:"week_day" validate:"min=0,max=6"`
	From    string `json:"from" validate:"required"`
	To      string `json:"to" validate:"required"`
	AllDay  bool   `json:"all_day"`
	Comment string `json:"comment"`
}

type updateScheduleRequest struct {
	Schedule []scheduleRequestDTO `json:"schedule" validate:"required,len=7,dive"`
}

func (r *updateScheduleRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}

	if len(r.Schedule) != 7 {
		return domainErr.NewValidationError("schedule must contain exactly 7 days").
			WithDetail("length", len(r.Schedule))
	}

	return nil
}

func (r *updateScheduleRequest) ToDomain(ses *session.Session) (*entity.User, error) {
	schedules := make([]entity.StaffSchedule, 0, len(r.Schedule))

	for _, s := range r.Schedule {
		opens, err := timeutil.ConvertAlmatyTimeToUTC(s.From)
		if err != nil {
			return nil, domainErr.NewValidationError("invalid time format in schedule").
				WithDetail("from", s.From).
				WithError(err)
		}

		closes, err := timeutil.ConvertAlmatyTimeToUTC(s.To)
		if err != nil {
			return nil, domainErr.NewValidationError("invalid time format in schedule").
				WithDetail("to", s.To).
				WithError(err)
		}

		schedules = append(schedules, entity.StaffSchedule{
			WeekDay: s.WeekDay,
			Open:    opens,
			Close:   closes,
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}

	return &entity.User{
		Phone:     ses.Phone(),
		Schedule:  schedules,
		UpdatedBy: ses.Phone(),
	}, nil
}
